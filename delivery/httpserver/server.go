package httpserver

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.project/go-fundamentals/gameapp/config/httpservercfg"
	"golang.project/go-fundamentals/gameapp/delivery/httpserver/userhandler"
	"log/slog"
	"net/http"
)

type HttpServer struct {
	ServerConfig httpservercfg.Config
	UserHandler  userhandler.UserHandler
}

func NewHttpServer(cfg httpservercfg.Config, userHandler userhandler.UserHandler) *HttpServer {

	return &HttpServer{ServerConfig: cfg, UserHandler: userHandler}
}

func (hs *HttpServer) Serve() {

	// Echo instance
	e := echo.New()

	// Middleware
	//e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/health-check", hs.HealthCheckHandler, middleware.Logger())

	hs.UserHandler.SetRoute(e)

	serverAddress := fmt.Sprintf("%s:%d", hs.ServerConfig.Host, hs.ServerConfig.Port)
	if err := e.Start(serverAddress); err != nil && errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server", "error", err)
	}
}
