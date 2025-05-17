package userhandler

import (
	"github.com/labstack/echo/v4"
	"golang.project/go-fundamentals/gameapp/config/httpservercfg/constant"
	"golang.project/go-fundamentals/gameapp/delivery/httpserver/parsericherror"
	"golang.project/go-fundamentals/gameapp/param"
	"golang.project/go-fundamentals/gameapp/service/authentication"
	"net/http"
)

func getClaims(ctx echo.Context) *authentication.Claims {

	claims, ok := ctx.Get(constant.AuthMiddlewareContextKey).(*authentication.Claims)
	if !ok {

		panic("JWT token missing or invalid")
	}

	return claims
}

func (h *UserHandler) userProfileHandler(ctx echo.Context) error {

	claims := getClaims(ctx)

	profile, profileErr := h.userService.Profile(param.NewProfileRequest(claims.UserId))
	if profileErr != nil {

		parseRichErr := parsericherror.New()
		message, statusCode := parseRichErr.ParseRichError(profileErr)

		return echo.NewHTTPError(statusCode, message)
	}

	return ctx.JSON(http.StatusFound, profile)
}
