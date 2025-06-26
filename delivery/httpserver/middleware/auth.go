package middleware

import (
	"errors"
	jwtMiddleware "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"golang.project/go-fundamentals/gameapp/config/httpservercfg/constant"
	"golang.project/go-fundamentals/gameapp/pkg/errormessage"
)

func (m *Middleware) AuthMiddleware() echo.MiddlewareFunc {
	return jwtMiddleware.WithConfig(
		jwtMiddleware.Config{
			ContextKey:    constant.AuthMiddlewareContextKey,
			SigningKey:    m.authService.Jwt.Config.SignKey,
			SigningMethod: m.authService.Jwt.Config.SignMethod,
			ParseTokenFunc: func(c echo.Context, authHeader string) (interface{}, error) {
				jwtToken := m.authService.Jwt.ExtractTokenFromHeader(authHeader)
				if jwtToken == "" {
					return nil, errors.New(errormessage.ErrorMsgEmptyJWT)
				}

				claims, err := m.authService.Jwt.ParseJWT(jwtToken)
				if err != nil {
					return nil, err
				}

				return claims, nil
			},
		},
	)
}
