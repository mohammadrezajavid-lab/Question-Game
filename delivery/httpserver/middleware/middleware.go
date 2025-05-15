package middleware

import "golang.project/go-fundamentals/gameapp/service/auth"

type Middleware struct {
	authService *auth.Service
}

func NewMiddleware(authService *auth.Service) Middleware {

	return Middleware{
		authService: authService,
	}
}
