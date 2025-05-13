package userhandler

import (
	"github.com/labstack/echo/v4"
	"golang.project/go-fundamentals/gameapp/delivery/httpserver/parsericherror"
	"golang.project/go-fundamentals/gameapp/dto"
	"golang.project/go-fundamentals/gameapp/pkg/normalize"
	"net/http"
)

func (h *UserHandler) userLoginHandler(ctx echo.Context) error {

	var requestUser = dto.NewLoginRequest()
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
	if validateErr, fieldErrors := h.UserValidator.ValidateLoginRequest(requestUser); validateErr != nil {

		parseErr := parsericherror.New()
		message, statusCode := parseErr.ParseRichError(validateErr)

		return ctx.JSON(statusCode, echo.Map{
			"message": message,
			"errors":  fieldErrors,
		})
	}

	loginRes, loginErr := h.UserService.Login(requestUser)
	if loginErr != nil {

		parseRichErr := parsericherror.New()
		message, statusCode := parseRichErr.ParseRichError(loginErr)

		return echo.NewHTTPError(statusCode, message)
	}

	return ctx.JSON(http.StatusOK, loginRes)
}
