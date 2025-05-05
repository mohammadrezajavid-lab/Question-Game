package httpserver

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gocasts.ir/go-fundamentals/gameapp/config"
	"gocasts.ir/go-fundamentals/gameapp/service/authorize"
	"gocasts.ir/go-fundamentals/gameapp/service/user"
	"log/slog"
	"net/http"
)

type HttpServer struct {
	Config      config.Config
	UserService *user.Service
	AuthService *authorize.Service
}

func NewHttpServer(cfg config.Config, userSvc *user.Service, authSvc *authorize.Service) *HttpServer {

	return &HttpServer{Config: cfg, UserService: userSvc, AuthService: authSvc}
}

func (hs *HttpServer) Serve() {

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/health-check", hs.HealthCheckHandler)
	e.POST("/users/register", hs.UserRegisterHandler)

	if err := e.Start(fmt.Sprintf("%s:%d", hs.Config.HttpServerCfg.Host, hs.Config.HttpServerCfg.Port)); err != nil && errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server", "error", err)
	}
}
