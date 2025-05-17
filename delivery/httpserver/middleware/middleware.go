package middleware

import "golang.project/go-fundamentals/gameapp/service/authentication"

type Middleware struct {
	authService *authentication.Service
}

func NewMiddleware(authService *authentication.Service) Middleware {

	return Middleware{
		authService: authService,
	}
}
