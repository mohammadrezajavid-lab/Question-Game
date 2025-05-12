package mysql

import (
	"database/sql"
	"errors"
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/pkg/richerror"
	"time"
)

func (d *DB) IsPhoneNumberUniq(phoneNumber string) (bool, error) {

	const operation = "mysql.IsPhoneNumberUniq"

	userRow := d.MysqlConnection.QueryRow(
		`SELECT * FROM game_app_db.Users WHERE phone_number = ?`,
		phoneNumber,
	)

	_, err := scanUser(userRow)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			return true, nil
		}

		return false, richerror.NewRichError(operation).
			WithError(err).
			WithMessage("unexpected error").
			WithKind(richerror.KindUnexpected)
	}

	return false, nil
}

func (d *DB) RegisterUser(user *entity.User) (*entity.User, error) {

	const operation = "mysql.RegisterUser"
	var result, eErr = d.MysqlConnection.Exec(
		`INSERT INTO game_app_db.Users(name, phone_number, hashed_password) VALUES(?, ?, ?)`,
		user.Name,
		user.PhoneNumber,
		user.HashedPassword,
	)
	if eErr != nil {

		return nil, richerror.NewRichError(operation).
			WithError(eErr).
			WithMessage("unexpected error: can't execute command").
			WithKind(richerror.KindUnexpected)
	}

	// error is always nil
	var userId, _ = result.LastInsertId()
	user.ID = uint(userId)

	return user, nil
}

func (d *DB) GetUserByPhoneNumber(phoneNumber string) (*entity.User, error) {

	const operation = "mysql.GetUserByPhoneNumber"

	var userRow = d.MysqlConnection.QueryRow(
		`SELECT * FROM game_app_db.Users WHERE phone_number = ?`,
		phoneNumber,
	)

	user, err := scanUser(userRow)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			// in this case, user record in database not found
			return nil, richerror.NewRichError(operation).
				WithError(err).
				WithMessage("record not found").
				WithKind(richerror.KindNotFound)
		}

		// in this case can't scan query from database
		return nil, richerror.NewRichError(operation).
			WithError(err).
			WithMessage("unexpected error: can't scan query result").
			WithKind(richerror.KindUnexpected)
	}

	return user, nil
}

func (d *DB) GetUserById(userId uint) (*entity.User, error) {

	const operation = "mysql.GetUserById"
	userRow := d.MysqlConnection.QueryRow(
		`SELECT * FROM game_app_db.Users WHERE id=?`,
		userId,
	)

	user, err := scanUser(userRow)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			return nil, richerror.NewRichError(operation).
				WithError(err).
				WithMessage("record not found").
				WithKind(richerror.KindNotFound).
				WithMeta(map[string]interface{}{"user_id": userId})
		}

		return nil, richerror.NewRichError(operation).
			WithError(err).
			WithMessage("unexpected error: can't scan query result").
			WithKind(richerror.KindUnexpected)
	}

	return user, nil
}

func scanUser(row *sql.Row) (*entity.User, error) {

	var createdAt time.Time

	user := entity.NewUser("", "", "")
	err := row.Scan(&user.ID, &user.Name, &user.PhoneNumber, &user.HashedPassword, &createdAt)

	return user, err
}
