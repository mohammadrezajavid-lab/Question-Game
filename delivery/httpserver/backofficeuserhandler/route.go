package backofficeuserhandler

import (
	"github.com/labstack/echo/v4"
	"golang.project/go-fundamentals/gameapp/delivery/httpserver/middleware"
	"golang.project/go-fundamentals/gameapp/entity"
)

func (h *BackOfficeUserHandler) SetRoute(e *echo.Echo) {

	newMiddleware := middleware.NewMiddleware(h.authService, h.authorizationService, h.presenceClient)

	userGroup := e.Group("/backoffice/users")

	userGroup.GET("/", h.listUsers, newMiddleware.AuthMiddleware(), newMiddleware.AccessCheck(entity.UserListPermission))
}
