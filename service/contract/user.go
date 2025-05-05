package contract

import "gocasts.ir/go-fundamentals/gameapp/entity"

type UserRepository interface {
	IsPhoneNumberUniq(phoneNumber string) (bool, error)
	RegisterUser(user *entity.User) (*entity.User, error)
	GetUserByPhoneNumber(phoneNumber string) (*entity.User, error)
	GetUserById(userId uint) (*entity.User, error)
}

type AuthorizeGenerator interface {
	CreateAccessToken(user *entity.User) (string, error)
	CreateRefreshToken(user *entity.User) (string, error)
}
