package matchinghandler

import (
	"github.com/labstack/echo/v4"
	"golang.project/go-fundamentals/gameapp/delivery/httpserver/middleware"
)

func (h *MatchingHandler) SetRoute(e *echo.Echo) {

	newMiddleware := middleware.NewMiddleware(h.authService, h.authorizationService, h.presenceClient)

	userGroup := e.Group("/matching-player/")

	userGroup.POST("add-to-waiting-list", h.addToWaitingList, newMiddleware.AuthMiddleware())
}
