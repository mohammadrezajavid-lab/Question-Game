package userhandler

import (
	"golang.project/go-fundamentals/gameapp/service/authentication"
	"golang.project/go-fundamentals/gameapp/service/user"
	"golang.project/go-fundamentals/gameapp/validator/uservalidator"
)

type UserHandler struct {
	userService *user.Service
	authService *authentication.Service

	userValidator *uservalidator.Validator
}

func NewHandler(userSvc *user.Service, authSvc *authentication.Service, userValidator *uservalidator.Validator) UserHandler {
	return UserHandler{
		userService:   userSvc,
		authService:   authSvc,
		userValidator: userValidator,
	}
}
