package userhandler

import (
	"golang.project/go-fundamentals/gameapp/service/auth"
	"golang.project/go-fundamentals/gameapp/service/user"
	"golang.project/go-fundamentals/gameapp/validator/uservalidator"
)

type UserHandler struct {
	userService *user.Service
	authService *auth.Service

	userValidator *uservalidator.Validator
}

func NewHandler(userSvc *user.Service, authSvc *auth.Service, userValidator *uservalidator.Validator) UserHandler {
	return UserHandler{
		userService:   userSvc,
		authService:   authSvc,
		userValidator: userValidator,
	}
}
