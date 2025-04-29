package mysql

import (
	"gocasts.ir/go-fundamentals/gameapp/entity"
)

func (d *DB) IsPhoneNumberUniq(phoneNumber string) (bool, error) {

	return true, nil
}

func (d *DB) RegisterUser(user *entity.User) (*entity.User, error) {
	return nil, nil
}
