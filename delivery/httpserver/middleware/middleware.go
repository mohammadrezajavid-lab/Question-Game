package middleware

import (
	"golang.project/go-fundamentals/gameapp/adapter/presenceclient"
	"golang.project/go-fundamentals/gameapp/service/authenticationservice"
	"golang.project/go-fundamentals/gameapp/service/authorizationservice"
)

type Middleware struct {
	authService          authenticationservice.Service
	authorizationService authorizationservice.Service
	presenceClient       presenceclient.Client
}

func NewMiddleware(
	authService authenticationservice.Service,
	authorizationService authorizationservice.Service,
	presenceClient presenceclient.Client,
) Middleware {

	return Middleware{
		authService:          authService,
		authorizationService: authorizationService,
		presenceClient:       presenceClient,
	}
}
