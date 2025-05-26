package middleware

import (
	"golang.project/go-fundamentals/gameapp/service/authenticationservice"
	"golang.project/go-fundamentals/gameapp/service/authorizationservice"
	"golang.project/go-fundamentals/gameapp/service/presenceservice"
)

type Middleware struct {
	authService          *authenticationservice.Service
	authorizationService *authorizationservice.Service
	presenceService      *presenceservice.Service
}

func NewMiddleware(
	authService *authenticationservice.Service,
	authorizationService *authorizationservice.Service,
	presenceService *presenceservice.Service,
) Middleware {

	return Middleware{
		authService:          authService,
		authorizationService: authorizationService,
		presenceService:      presenceService,
	}
}
