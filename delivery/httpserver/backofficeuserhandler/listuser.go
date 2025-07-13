package backofficeuserhandler

import (
	"github.com/labstack/echo/v4"
	"golang.project/go-fundamentals/gameapp/delivery/httpserver/parsericherror"
	"net/http"
)

func (h *BackOfficeUserHandler) listUsers(ctx echo.Context) error {

	users, lErr := h.backOfficeUserService.ListAllUsers(ctx.Request().Context())
	if lErr != nil {
		message, statusCode := parsericherror.New().ParseRichError(lErr)

		return ctx.JSON(statusCode, message)
	}

	return ctx.JSON(http.StatusOK, echo.Map{"data": users})
}
