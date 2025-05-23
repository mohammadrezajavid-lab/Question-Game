package userhandler

import (
	"github.com/labstack/echo/v4"
	"golang.project/go-fundamentals/gameapp/delivery/httpserver/parsericherror"
	"golang.project/go-fundamentals/gameapp/param"
	"golang.project/go-fundamentals/gameapp/pkg/normalize"
	"net/http"
)

func (h *UserHandler) userLoginHandler(ctx echo.Context) error {

	var requestUser = param.NewLoginRequest()
	if err := ctx.Bind(requestUser); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// normalized register request
	norm := normalize.New()
	phoneNumber, err := norm.NormalizePhoneNumber(requestUser.PhoneNumber)
	if err != nil {

		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}
	requestUser.PhoneNumber = phoneNumber

	// validate login request
	if validateErr, fieldErrors := h.userValidator.ValidateLoginRequest(requestUser); validateErr != nil {

		message, statusCode := parsericherror.New().ParseRichError(validateErr)

		return ctx.JSON(statusCode, echo.Map{
			"message": message,
			"errors":  fieldErrors,
		})
	}

	loginRes, loginErr := h.userService.Login(requestUser)
	if loginErr != nil {

		message, statusCode := parsericherror.New().ParseRichError(loginErr)

		return echo.NewHTTPError(statusCode, echo.Map{
			"message": message,
		})
	}

	return ctx.JSON(http.StatusOK, loginRes)
}
