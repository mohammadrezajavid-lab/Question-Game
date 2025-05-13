package userhandler

import (
	"github.com/labstack/echo/v4"
	"golang.project/go-fundamentals/gameapp/delivery/httpserver/parsericherror"
	"golang.project/go-fundamentals/gameapp/param"
	"golang.project/go-fundamentals/gameapp/pkg/constant"
	"golang.project/go-fundamentals/gameapp/service/auth"
	"net/http"
)

func getClaims(ctx echo.Context) *auth.Claims {

	claims, ok := ctx.Get(constant.AuthMiddlewareContextKey).(*auth.Claims)
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
