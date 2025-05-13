package userhandler

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (h *UserHandler) SetRoute(e *echo.Echo) {

	// Great one group and add logger middleware for userGroup
	userGroup := e.Group("/users/", middleware.Logger())

	userGroup.POST("register", h.userRegisterHandler)
	userGroup.POST("login", h.userLoginHandler)
	userGroup.GET("profile", h.userProfileHandler)
}
