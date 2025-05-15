package middleware

import (
	"github.com/golang-jwt/jwt/v4"
	jwtMiddleware "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"golang.project/go-fundamentals/gameapp/config/httpservercfg/constant"
)

func (m *Middleware) AuthMiddleware() echo.MiddlewareFunc {
	return jwtMiddleware.WithConfig(
		jwtMiddleware.Config{
			ContextKey:    constant.AuthMiddlewareContextKey,
			SigningKey:    m.authService.Config.SignKey,
			SigningMethod: jwt.SigningMethodES256.Alg(),
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
