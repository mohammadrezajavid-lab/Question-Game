package claim

import (
	"errors"
	"github.com/labstack/echo/v4"
	"golang.project/go-fundamentals/gameapp/config/httpservercfg/constant"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/pkg/errormessage"
	"golang.project/go-fundamentals/gameapp/service/authenticationservice"
)

func GetClaimsFromEchoContext(ctx echo.Context) *authenticationservice.Claims {

	claims, ok := ctx.Get(constant.AuthMiddlewareContextKey).(*authenticationservice.Claims)
	if !ok {
		logger.Panic(errors.New(errormessage.ErrorMsgInvalidJWT), errormessage.ErrorMsgInvalidJWT)
	}

	return claims
}
