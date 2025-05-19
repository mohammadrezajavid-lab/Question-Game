package claim

import (
	"github.com/labstack/echo/v4"
	"golang.project/go-fundamentals/gameapp/config/httpservercfg/constant"
	"golang.project/go-fundamentals/gameapp/service/authenticationservice"
)

func GetClaimsFromEchoContext(ctx echo.Context) *authenticationservice.Claims {

	claims, ok := ctx.Get(constant.AuthMiddlewareContextKey).(*authenticationservice.Claims)
	if !ok {

		panic("JWT token missing or invalid")
	}

	return claims
}
