package middleware

import (
	jwtMiddleware "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"golang.project/go-fundamentals/gameapp/config/httpservercfg/constant"
)

func (m *Middleware) AuthMiddleware() echo.MiddlewareFunc {
	return jwtMiddleware.WithConfig(
		jwtMiddleware.Config{
			ContextKey:    constant.AuthMiddlewareContextKey,
			SigningKey:    m.authService.Config.SignKey,
			SigningMethod: constant.DefaultSignMethod,
			ParseTokenFunc: func(c echo.Context, auth string) (interface{}, error) {
				claims, err := m.authService.ParseJWT(auth)
				if err != nil {
					return nil, err
				}

				return claims, nil
			},
		},
	)
}
