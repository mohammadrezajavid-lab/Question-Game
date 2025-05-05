package httpserver

import (
	"github.com/labstack/echo/v4"
	"gocasts.ir/go-fundamentals/gameapp/service/user"
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
