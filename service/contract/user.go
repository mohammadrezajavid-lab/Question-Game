package contract

import "gocasts.ir/go-fundamentals/gameapp/entity"

type UserRepository interface {
	IsPhoneNumberUniq(phoneNumber string) (bool, error)
	CreateNewUser(user *entity.User) (*entity.User, error)
}
