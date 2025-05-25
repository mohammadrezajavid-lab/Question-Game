package userservice

import (
	"context"
	"golang.project/go-fundamentals/gameapp/entity"
)

type Repository interface {
	RegisterUser(user *entity.User) (*entity.User, error)
	GetUserByPhoneNumber(phoneNumber string) (*entity.User, error)
	GetUserById(ctx context.Context, userId uint) (*entity.User, error)
}

type AuthGenerator interface {
	CreateAccessToken(user *entity.User) (string, error)
	CreateRefreshToken(user *entity.User) (string, error)
}

type Service struct {
	userRepository Repository
	authService    AuthGenerator
}

func NewService(userRepository Repository, authService AuthGenerator) *Service {
	return &Service{
		userRepository: userRepository,
		authService:    authService,
	}
}

// All Request Inputs for Interaction/Service Should be Sanitized.
