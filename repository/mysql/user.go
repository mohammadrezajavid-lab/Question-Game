package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"golang.project/go-fundamentals/gameapp/entity"
	"time"
)

func (d *DB) IsPhoneNumberUniq(phoneNumber string) (bool, error) {

	userRow := d.MysqlConnection.QueryRow(
		`SELECT * FROM game_app_db.Users WHERE phone_number = ?`,
		phoneNumber,
	)

	_, err := scanUser(userRow)
	if errors.Is(err, sql.ErrNoRows) {

		return true, nil
	}

	return false, nil
}

func (d *DB) RegisterUser(user *entity.User) (*entity.User, error) {

	var result, eErr = d.MysqlConnection.Exec(
		`INSERT INTO game_app_db.Users(name, phone_number, hashed_password) VALUES(?, ?, ?)`,
		user.Name,
		user.PhoneNumber,
		user.HashedPassword,
	)
	if eErr != nil {

		return nil, fmt.Errorf("can't execute command: %w", eErr)
	}

	// error is always nil
	var userId, _ = result.LastInsertId()
	user.ID = uint(userId)

	return user, nil
}

func (d *DB) GetUserByPhoneNumber(phoneNumber string) (*entity.User, error) {

	var userRow = d.MysqlConnection.QueryRow(
		`SELECT * FROM game_app_db.Users WHERE phone_number = ?`,
		phoneNumber,
	)

	user, err := scanUser(userRow)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			return nil, fmt.Errorf("record not found\n")
		}

		return nil, fmt.Errorf("can't scan query result: %w\n", err)
	}

	return user, nil
}

func (d *DB) GetUserById(userId uint) (*entity.User, error) {

	userRow := d.MysqlConnection.QueryRow(
		`SELECT * FROM game_app_db.Users WHERE id=?`,
		userId,
	)

	user, err := scanUser(userRow)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			return nil, fmt.Errorf("record not found")
		}

		return nil, fmt.Errorf("can't scan query result: %w", err)
	}

	return user, nil
}

func scanUser(row *sql.Row) (*entity.User, error) {

	var createdAt time.Time

	user := entity.NewUser("", "", "")
	err := row.Scan(&user.ID, &user.Name, &user.PhoneNumber, &user.HashedPassword, &createdAt)

	return user, err
}
