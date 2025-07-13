package usermysql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/prometheus/client_golang/prometheus"
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/metrics"
	"golang.project/go-fundamentals/gameapp/pkg/errormessage"
	"golang.project/go-fundamentals/gameapp/pkg/richerror"
	"golang.project/go-fundamentals/gameapp/repository/mysql"
	"time"
)

func (d *DataBase) IsPhoneNumberUniq(phoneNumber string) (bool, error) {

	const operation = "mysql.user.IsPhoneNumberUniq"
	const queryType = "select"

	userRow := d.dataBase.MysqlConnection.QueryRow(
		`SELECT * FROM game_app_db.users WHERE phone_number = ?`,
		phoneNumber,
	)

	metrics.DBQueryCounter.With(prometheus.Labels{"query_type": queryType}).Inc()

	_, err := scanUser(userRow)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return true, nil
		}

		metrics.DBFailedQueryCounter.With(prometheus.Labels{"query_type": queryType}).Inc()

		return false, richerror.NewRichError(operation).
			WithError(err).
			WithMessage(errormessage.ErrorMsgUnexpected).
			WithKind(richerror.KindUnexpected)
	}

	return false, nil
}

func (d *DataBase) RegisterUser(user *entity.User) (*entity.User, error) {

	const operation = "mysql.user.RegisterUser"
	const queryType = "insert"

	var result, eErr = d.dataBase.MysqlConnection.Exec(
		`INSERT INTO game_app_db.users(name, phone_number, hashed_password, role) VALUES(?, ?, ?, ?)`,
		user.Name,
		user.PhoneNumber,
		user.HashedPassword,
		user.Role.String(),
	)

	metrics.DBQueryCounter.With(prometheus.Labels{"query_type": queryType}).Inc()

	if eErr != nil {
		logger.Warn(eErr, errormessage.ErrorMsgFailedExecuteQuery)

		metrics.DBFailedQueryCounter.With(prometheus.Labels{"query_type": queryType}).Inc()

		return nil, richerror.NewRichError(operation).
			WithError(eErr).
			WithMessage(errormessage.ErrorMsgUnexpected).
			WithKind(richerror.KindUnexpected)
	}

	// error is always nil
	var userId, _ = result.LastInsertId()
	user.Id = uint(userId)

	return user, nil
}

func (d *DataBase) GetUserByPhoneNumber(phoneNumber string) (*entity.User, error) {

	const operation = "mysql.user.GetUserByPhoneNumber"
	const queryType = "select"

	var userRow = d.dataBase.MysqlConnection.QueryRow(
		`SELECT * FROM game_app_db.users WHERE phone_number = ?`,
		phoneNumber,
	)

	metrics.DBQueryCounter.With(prometheus.Labels{"query_type": queryType}).Inc()

	user, err := scanUser(userRow)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			// in this case, user record in database not found
			return nil, richerror.NewRichError(operation).
				WithError(err).
				WithMessage(errormessage.ErrorMsgRecordNotFound).
				WithKind(richerror.KindNotFound)
		}

		// in this case can't scan query from database
		metrics.DBFailedQueryCounter.With(prometheus.Labels{"query_type": queryType}).Inc()

		return nil, richerror.NewRichError(operation).
			WithError(err).
			WithMessage(errormessage.ErrorMsgUnexpected).
			WithKind(richerror.KindUnexpected)
	}

	return user, nil
}

func (d *DataBase) GetUserById(ctx context.Context, userId uint) (*entity.User, error) {

	const operation = "mysql.user.GetUserById"
	const queryType = "select"

	userRow := d.dataBase.MysqlConnection.QueryRowContext(
		ctx,
		`SELECT * FROM game_app_db.users WHERE id=?`,
		userId,
	)

	metrics.DBQueryCounter.With(prometheus.Labels{"query_type": queryType}).Inc()

	user, err := scanUser(userRow)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			return nil, richerror.NewRichError(operation).
				WithError(err).
				WithMessage(errormessage.ErrorMsgRecordNotFound).
				WithKind(richerror.KindNotFound).
				WithMeta(map[string]interface{}{"user_id": userId})
		}

		metrics.DBFailedQueryCounter.With(prometheus.Labels{"query_type": queryType}).Inc()

		return nil, richerror.NewRichError(operation).
			WithError(err).
			WithMessage(errormessage.ErrorMsgUnexpected).
			WithKind(richerror.KindUnexpected)
	}

	return user, nil
}

func (d *DataBase) ListUsers(ctx context.Context) ([]entity.User, error) {
	const operation = "mysql.user.ListUsers"
	const queryType = "select"

	userRows, err := d.dataBase.MysqlConnection.QueryContext(ctx, `SELECT * FROM game_app_db.users`)
	metrics.DBQueryCounter.With(prometheus.Labels{"query_type": queryType}).Inc()
	if err != nil {
		metrics.DBFailedQueryCounter.With(prometheus.Labels{"query_type": queryType}).Inc()
		return nil, richerror.NewRichError(operation).
			WithError(err).
			WithMessage(errormessage.ErrorMsgUnexpected).
			WithKind(richerror.KindUnexpected)
	}
	defer userRows.Close()

	users := make([]entity.User, 0)

	for userRows.Next() {
		user, sErr := scanUser(userRows)
		if sErr != nil {
			if errors.Is(sErr, sql.ErrNoRows) {

				return nil, richerror.NewRichError(operation).
					WithError(sErr).
					WithMessage(errormessage.ErrorMsgRecordNotFound).
					WithKind(richerror.KindNotFound).
					WithMeta(map[string]interface{}{"user_id": user.Id})
			}

			metrics.DBFailedQueryCounter.With(prometheus.Labels{"query_type": queryType}).Inc()

			return nil, richerror.NewRichError(operation).
				WithError(sErr).
				WithMessage(errormessage.ErrorMsgUnexpected).
				WithKind(richerror.KindUnexpected)
		}

		users = append(users, *user)
	}

	if rErr := userRows.Err(); rErr != nil {
		metrics.DBFailedQueryCounter.With(prometheus.Labels{"query_type": queryType}).Inc()
		return nil, richerror.NewRichError(operation).
			WithError(rErr).
			WithMessage(rErr.Error())
	}

	return users, nil
}

func scanUser(scanner mysql.Scanner) (*entity.User, error) {

	var createdAt time.Time
	var roleStr string

	user := entity.NewUser("", "", "")
	err := scanner.Scan(&user.Id, &user.Name, &user.PhoneNumber, &user.HashedPassword, &createdAt, &roleStr)

	if err == nil {
		role := entity.MapToRoleEntity(roleStr)
		if role == 0 {

			return user, errors.New(errormessage.ErrorMsgScanQuery)
		}
		user.Role = role
	}

	return user, err
}
