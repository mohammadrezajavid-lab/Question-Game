package userhandler

import (
	"golang.project/go-fundamentals/gameapp/service/auth"
	"golang.project/go-fundamentals/gameapp/service/user"
	"golang.project/go-fundamentals/gameapp/validator/uservalidator"
)

type UserHandler struct {
	UserService *user.Service
	AuthService *auth.Service

	UserValidator *uservalidator.Validator
}

func NewHandler(userSvc *user.Service, authSvc *auth.Service, userValidator *uservalidator.Validator) UserHandler {
	return UserHandler{
		UserService:   userSvc,
		AuthService:   authSvc,
		UserValidator: userValidator,
	}
}
