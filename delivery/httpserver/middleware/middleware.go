package middleware

import (
	"golang.project/go-fundamentals/gameapp/service/authenticationservice"
	"golang.project/go-fundamentals/gameapp/service/authorizationservice"
)

type Middleware struct {
	authService          *authenticationservice.Service
	authorizationService *authorizationservice.Service
}

func NewMiddleware(
	authService *authenticationservice.Service,
	authorizationService *authorizationservice.Service,
) Middleware {

	return Middleware{
		authService:          authService,
		authorizationService: authorizationService,
	}
}
