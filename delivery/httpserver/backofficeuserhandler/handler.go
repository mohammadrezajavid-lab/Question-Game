package backofficeuserhandler

import (
	"golang.project/go-fundamentals/gameapp/adapter/presenceclient"
	"golang.project/go-fundamentals/gameapp/service/authenticationservice"
	"golang.project/go-fundamentals/gameapp/service/authorizationservice"
	"golang.project/go-fundamentals/gameapp/service/backofficeuserservice"
	"golang.project/go-fundamentals/gameapp/validator/uservalidator"
)

type BackOfficeUserHandler struct {
	backOfficeUserService backofficeuserservice.Service
	authService           authenticationservice.Service
	authorizationService  authorizationservice.Service
	userValidator         uservalidator.Validator
	presenceClient        presenceclient.Client
}

func NewHandler(
	backOfficeUserSvc backofficeuserservice.Service,
	authSvc authenticationservice.Service,
	authorizationSvc authorizationservice.Service,
	userValidator uservalidator.Validator,
	presenceClient presenceclient.Client,
) BackOfficeUserHandler {
	return BackOfficeUserHandler{
		backOfficeUserService: backOfficeUserSvc,
		authService:           authSvc,
		authorizationService:  authorizationSvc,
		userValidator:         userValidator,
		presenceClient:        presenceClient,
	}
}
