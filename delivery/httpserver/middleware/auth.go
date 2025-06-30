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
			SigningKey:    m.authService.Jwt.Config.SignKey,
			SigningMethod: m.authService.Jwt.Config.SignMethod,
			ParseTokenFunc: func(c echo.Context, jwtToken string) (interface{}, error) {
				claims, err := m.authService.Jwt.ParseJWT(jwtToken)
				if err != nil {
					return nil, err
				}
				return claims, nil
			},
		},
	)
}
