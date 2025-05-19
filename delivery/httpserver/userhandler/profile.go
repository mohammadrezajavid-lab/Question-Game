package userhandler

import (
	"github.com/labstack/echo/v4"
	"golang.project/go-fundamentals/gameapp/delivery/httpserver/parsericherror"
	"golang.project/go-fundamentals/gameapp/param"
	"golang.project/go-fundamentals/gameapp/pkg/claim"
	"net/http"
)

func (h *UserHandler) userProfileHandler(ctx echo.Context) error {

	claims := claim.GetClaimsFromEchoContext(ctx)

	profile, profileErr := h.userService.Profile(param.NewProfileRequest(claims.UserId))
	if profileErr != nil {

		message, statusCode := parsericherror.New().ParseRichError(profileErr)

		return echo.NewHTTPError(statusCode, message)
	}

	return ctx.JSON(http.StatusFound, profile)
}
