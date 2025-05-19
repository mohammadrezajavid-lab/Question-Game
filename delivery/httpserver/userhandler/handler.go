package userhandler

import (
	"golang.project/go-fundamentals/gameapp/service/authenticationservice"
	"golang.project/go-fundamentals/gameapp/service/authorizationservice"
	"golang.project/go-fundamentals/gameapp/service/userservice"
	"golang.project/go-fundamentals/gameapp/validator/uservalidator"
)

type UserHandler struct {
	userService          *userservice.Service
	authService          *authenticationservice.Service
	authorizationService *authorizationservice.Service

	userValidator *uservalidator.Validator
}

func NewHandler(
	userSvc *userservice.Service,
	authSvc *authenticationservice.Service,
	authorizationSvc *authorizationservice.Service,
	userValidator *uservalidator.Validator,
) UserHandler {
	return UserHandler{
		userService:          userSvc,
		authService:          authSvc,
		authorizationService: authorizationSvc,
		userValidator:        userValidator,
	}
}
