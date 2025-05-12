package httpserver

import (
	"github.com/labstack/echo/v4"
	"golang.project/go-fundamentals/gameapp/delivery/httpserver/parsericherror"
	"golang.project/go-fundamentals/gameapp/dto"
	"golang.project/go-fundamentals/gameapp/pkg/normalize"
	"net/http"
)

type Response struct {
	Message string           `json:"message"`
	Errors  map[string]error `json:"errors"`
}

func (hs *HttpServer) UserRegisterHandler(ctx echo.Context) error {

	var requestUser = dto.NewRegisterRequest()
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

	// validate register request
	if validateErr, fieldErrors := hs.UserValidator.ValidateRegisterRequest(requestUser); validateErr != nil {

		parseRichErr := parsericherror.New()
		message, statusCode := parseRichErr.ParseRichError(validateErr)

		return ctx.JSON(statusCode, echo.Map{
			"message": message,
			"errors":  fieldErrors,
		})
	}

	registerResponse, registerErr := hs.UserService.Register(requestUser)
	if registerErr != nil {

		parseRichErr := parsericherror.New()
		message, statusCode := parseRichErr.ParseRichError(registerErr)

		return echo.NewHTTPError(statusCode, message)
	}

	return ctx.JSON(http.StatusCreated, registerResponse)
}

func (hs *HttpServer) UserLoginHandler(ctx echo.Context) error {

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
	if validateErr, fieldErrors := hs.UserValidator.ValidateLoginRequest(requestUser); validateErr != nil {

		parseErr := parsericherror.New()
		message, statusCode := parseErr.ParseRichError(validateErr)

		return ctx.JSON(statusCode, echo.Map{
			"message": message,
			"errors":  fieldErrors,
		})
	}

	loginRes, loginErr := hs.UserService.Login(requestUser)
	if loginErr != nil {

		parseRichErr := parsericherror.New()
		message, statusCode := parseRichErr.ParseRichError(loginErr)

		return echo.NewHTTPError(statusCode, message)
	}

	return ctx.JSON(http.StatusOK, loginRes)
}

func (hs *HttpServer) UserProfileHandler(ctx echo.Context) error {

	// TODO - we are sanitize userId in this handler after send userId to service layer

	req := ctx.Request()
	tokenAuth := req.Header.Get("Authorization")
	claims, parseJWTErr := hs.AuthService.ParseJWT(tokenAuth)
	if parseJWTErr != nil {

		return echo.NewHTTPError(http.StatusUnauthorized, parseJWTErr.Error())
	}

	if claims == nil {

		return echo.NewHTTPError(http.StatusUnauthorized, "claims is empty")
	}

	profile, profileErr := hs.UserService.Profile(dto.NewProfileRequest(claims.UserId))
	if profileErr != nil {

		parseRichErr := parsericherror.New()
		message, statusCode := parseRichErr.ParseRichError(profileErr)

		return echo.NewHTTPError(statusCode, message)
	}

	return ctx.JSON(http.StatusFound, profile)
}
