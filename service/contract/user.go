package contract

import (
	"golang.project/go-fundamentals/gameapp/entity"
)

type UserRepository interface {
	RegisterUser(user *entity.User) (*entity.User, error)
	GetUserByPhoneNumber(phoneNumber string) (*entity.User, bool, error)
	GetUserById(userId uint) (*entity.User, error)
}

type AuthorizeGenerator interface {
	CreateAccessToken(user *entity.User) (string, error)
	CreateRefreshToken(user *entity.User) (string, error)
}
