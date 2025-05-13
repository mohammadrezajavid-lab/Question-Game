package userhandler

import (
	"github.com/labstack/echo/v4"
	"golang.project/go-fundamentals/gameapp/delivery/httpserver/parsericherror"
	"golang.project/go-fundamentals/gameapp/dto"
	"net/http"
)

func (h *UserHandler) userProfileHandler(ctx echo.Context) error {

	// TODO - we are sanitize userId in this handler after send userId to service layer

	req := ctx.Request()
	tokenAuth := req.Header.Get("Authorization")
	claims, parseJWTErr := h.AuthService.ParseJWT(tokenAuth)
	if parseJWTErr != nil {

		return echo.NewHTTPError(http.StatusUnauthorized, parseJWTErr.Error())
	}

	if claims == nil {

		return echo.NewHTTPError(http.StatusUnauthorized, "claims is empty")
	}

	profile, profileErr := h.UserService.Profile(dto.NewProfileRequest(claims.UserId))
	if profileErr != nil {

		parseRichErr := parsericherror.New()
		message, statusCode := parseRichErr.ParseRichError(profileErr)

		return echo.NewHTTPError(statusCode, message)
	}

	return ctx.JSON(http.StatusFound, profile)
}
