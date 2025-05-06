package httpserver

import (
	"github.com/labstack/echo/v4"
	"golang.project/go-fundamentals/gameapp/service/user"
	"net/http"
)

func (hs *HttpServer) UserRegisterHandler(ctx echo.Context) error {

	var requestUser = user.NewRegisterRequest()
	if err := ctx.Bind(requestUser); err != nil {

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	registerResponse, registerErr := hs.UserService.Register(requestUser)
	if registerErr != nil {

		return echo.NewHTTPError(http.StatusBadRequest, registerErr.Error())
	}

	return ctx.JSON(http.StatusCreated, registerResponse)
}

func (hs *HttpServer) UserLoginHandler(ctx echo.Context) error {

	var requestUser = user.NewLoginRequest("", "")
	if err := ctx.Bind(requestUser); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	loginRes, lErr := hs.UserService.Login(requestUser)
	if lErr != nil {

		return echo.NewHTTPError(http.StatusUnauthorized, lErr.Error())
	}

	return ctx.JSON(http.StatusOK, loginRes)
}

func (hs *HttpServer) UserProfileHandler(ctx echo.Context) error {

	// TODO - we are sanitize userId in this handler after send userId to service layer

	req := ctx.Request()
	tokenAuth := req.Header.Get("Authorization")
	claims, parseErr := hs.AuthService.ParseJWT(tokenAuth)
	if parseErr != nil {

		return echo.NewHTTPError(http.StatusUnauthorized, parseErr.Error())
	}

	if claims == nil {

		return echo.NewHTTPError(http.StatusUnauthorized, "claims is empty")
	}

	profile, pErr := hs.UserService.Profile(user.NewProfileRequest(claims.UserId))
	if pErr != nil {

		return echo.NewHTTPError(http.StatusInternalServerError, pErr.Error())
	}

	return ctx.JSON(http.StatusFound, profile)
}
