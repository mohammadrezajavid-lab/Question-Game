package userhandler

import (
	"github.com/labstack/echo/v4"
	"golang.project/go-fundamentals/gameapp/delivery/httpserver/parsericherror"
	"golang.project/go-fundamentals/gameapp/param/userparam"
	"golang.project/go-fundamentals/gameapp/pkg/normalize"
	"net/http"
)

func (h *UserHandler) userRegisterHandler(ctx echo.Context) error {

	var requestUser = userparam.NewRegisterRequest()
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

	if fieldErrors, validateErr := h.userValidator.ValidateRegisterRequest(requestUser); validateErr != nil {

		message, statusCode := parsericherror.New().ParseRichError(validateErr)

		return ctx.JSON(statusCode, echo.Map{
			"message": message,
			"errors":  fieldErrors,
		})
	}

	registerResponse, registerErr := h.userService.Register(requestUser)
	if registerErr != nil {

		message, statusCode := parsericherror.New().ParseRichError(registerErr)

		return echo.NewHTTPError(statusCode, message)
	}

	return ctx.JSON(http.StatusCreated, registerResponse)
}
