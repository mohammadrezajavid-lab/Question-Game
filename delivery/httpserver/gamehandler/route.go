package gamehandler

import (
	"github.com/labstack/echo/v4"
	"golang.project/go-fundamentals/gameapp/delivery/httpserver/middleware"
)

func (q *GameHandler) SetRoute(e *echo.Echo) {
	newMiddleware := middleware.NewMiddleware(q.authService, q.authorizationService, q.presenceClient)

	gameGroup := e.Group("/game/")

	gameGroup.POST("add-to-waiting-list", q.addToWaitingList, newMiddleware.AuthMiddleware())
	gameGroup.POST("get-quiz", q.getQuiz, newMiddleware.AuthMiddleware())
}
