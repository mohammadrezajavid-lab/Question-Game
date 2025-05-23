package userhandler

import (
	"github.com/labstack/echo/v4"
	"golang.project/go-fundamentals/gameapp/delivery/httpserver/middleware"
)

func (h *UserHandler) SetRoute(e *echo.Echo) {

	newMiddleware := middleware.NewMiddleware(h.authService, h.authorizationService)

	userGroup := e.Group("/users/")

	userGroup.POST("register", h.userRegisterHandler)
	userGroup.POST("login", h.userLoginHandler)
	userGroup.GET("profile", h.userProfileHandler, newMiddleware.AuthMiddleware())
}
