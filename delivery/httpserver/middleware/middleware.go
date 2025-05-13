package middleware

import "golang.project/go-fundamentals/gameapp/service/auth"

type Middleware struct {
	authConfig  auth.Config
	authService *auth.Service
}

func NewMiddleware(authConfig auth.Config, authService *auth.Service) Middleware {
	return Middleware{
		authConfig:  authConfig,
		authService: authService,
	}
}
