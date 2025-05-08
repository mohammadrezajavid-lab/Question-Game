package httpserver

import (
	"github.com/labstack/echo/v4"
	"golang.project/go-fundamentals/gameapp/datatransferobject/userdto"
	"golang.project/go-fundamentals/gameapp/delivery/httpserver/parsericherror"
	"net/http"
)

func (hs *HttpServer) UserRegisterHandler(ctx echo.Context) error {

	var requestUser = userdto.NewRegisterRequest()
	if err := ctx.Bind(requestUser); err != nil {

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// validate requestUser
	err := hs.UserValidator.ValidateRegisterRequest(requestUser)
	if err != nil {
		parseRichErr := parsericherror.New()
		message, statusCode := parseRichErr.ParseRichError(err)

		return echo.NewHTTPError(statusCode, message)
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

	var requestUser = userdto.NewLoginRequest("", "")
	if err := ctx.Bind(requestUser); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	loginRes, lErr := hs.UserService.Login(requestUser)
	if lErr != nil {

		parseRichErr := parsericherror.New()
		message, statusCode := parseRichErr.ParseRichError(lErr)

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

	profile, pErr := hs.UserService.Profile(userdto.NewProfileRequest(claims.UserId))
	if pErr != nil {

		parseRichErr := parsericherror.New()
		message, statusCode := parseRichErr.ParseRichError(pErr)

		return echo.NewHTTPError(statusCode, message)
	}

	return ctx.JSON(http.StatusFound, profile)
}
