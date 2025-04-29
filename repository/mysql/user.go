package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"gocasts.ir/go-fundamentals/gameapp/entity"
)

func (d *DB) IsPhoneNumberUniq(phoneNumber string) (bool, error) {

	var rowUser = new(entity.User)
	var createdAt []byte

	var result *sql.Row = d.MysqlConnection.QueryRow(`SELECT * FROM game_app_db.Users WHERE phone_number = ?`, phoneNumber)

	if result.Err() != nil {
		return false, result.Err()
	}

	if scanErr := result.Scan(&rowUser.ID, &rowUser.Name, &rowUser.PhoneNumber, &createdAt); errors.Is(scanErr, sql.ErrNoRows) {
		return true, nil
	}

	return false, nil
}

func (d *DB) RegisterUser(user *entity.User) (*entity.User, error) {

	var result, eErr = d.MysqlConnection.Exec(`INSERT INTO game_app_db.Users(name, phone_number) VALUES(?, ?)`, user.Name, user.PhoneNumber)
	if eErr != nil {
		return entity.NewUser("", ""), fmt.Errorf("can't execute command: %w", eErr)
	}

	// error is always nil
	var userId, _ = result.LastInsertId()
	user.ID = uint(userId)

	return user, nil
}
