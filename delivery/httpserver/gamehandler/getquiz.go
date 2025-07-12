package gamehandler

import (
	"github.com/labstack/echo/v4"
	"golang.project/go-fundamentals/gameapp/delivery/httpserver/parsericherror"
	"golang.project/go-fundamentals/gameapp/param/gameparam"
	"net/http"
)

func (q *GameHandler) getQuiz(ctx echo.Context) error {

	var request = gameparam.GetQuizRequest{}
	if err := ctx.Bind(request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	responseQuiz, err := q.gameService.GetQuiz(ctx.Request().Context(), request)
	if err != nil {
		message, statusCode := parsericherror.New().ParseRichError(err)
		return echo.NewHTTPError(statusCode, echo.Map{
			"message": message,
		})
	}

	return ctx.JSON(http.StatusOK, responseQuiz)
}
