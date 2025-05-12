package user

import (
	"golang.project/go-fundamentals/gameapp/entity"
)

type Repository interface {
	RegisterUser(user *entity.User) (*entity.User, error)
	GetUserByPhoneNumber(phoneNumber string) (*entity.User, error)
	GetUserById(userId uint) (*entity.User, error)
}

type AuthorizeGenerator interface {
	CreateAccessToken(user *entity.User) (string, error)
	CreateRefreshToken(user *entity.User) (string, error)
}

type Service struct {
	userRepository Repository
	authService    AuthorizeGenerator
}

func NewService(userRepository Repository, authorizeService AuthorizeGenerator) *Service {
	return &Service{
		userRepository: userRepository,
		authService:    authorizeService,
	}
}

// All Request Inputs for Interaction/Service Should be Sanitized.
