package httpserver

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.project/go-fundamentals/gameapp/config"
	"golang.project/go-fundamentals/gameapp/service/auth"
	"golang.project/go-fundamentals/gameapp/service/user"
	"golang.project/go-fundamentals/gameapp/validator/uservalidator"
	"log/slog"
	"net/http"
)

type HttpServer struct {
	Config      config.Config
	UserService *user.Service
	AuthService *auth.Service

	UserValidator *uservalidator.Validator
}

func NewHttpServer(cfg config.Config, userSvc *user.Service, authSvc *auth.Service, userValidator *uservalidator.Validator) *HttpServer {

	return &HttpServer{Config: cfg, UserService: userSvc, AuthService: authSvc, UserValidator: userValidator}
}

func (hs *HttpServer) Serve() {

	// Echo instance
	e := echo.New()

	// Middleware
	//e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/health-check", hs.HealthCheckHandler)

	// Great one group and add logger middleware for userGroup
	userGroup := e.Group("/users/", middleware.Logger())
	userGroup.POST("register", hs.UserRegisterHandler)
	userGroup.POST("login", hs.UserLoginHandler)
	userGroup.GET("profile", hs.UserProfileHandler)

	serverAddress := fmt.Sprintf("%s:%d", hs.Config.HttpServerCfg.Host, hs.Config.HttpServerCfg.Port)
	if err := e.Start(serverAddress); err != nil && errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server", "error", err)
	}
}
