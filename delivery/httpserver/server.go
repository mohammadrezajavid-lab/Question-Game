package httpserver

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.project/go-fundamentals/gameapp/config/httpservercfg"
	"golang.project/go-fundamentals/gameapp/delivery/httpserver/userhandler"
	"golang.project/go-fundamentals/gameapp/service/authentication"
	"golang.project/go-fundamentals/gameapp/service/user"
	"golang.project/go-fundamentals/gameapp/validator/uservalidator"
	"log/slog"
	"net/http"
)

type HttpServer struct {
	serverConfig httpservercfg.Config
	userHandler  userhandler.UserHandler
}

func NewHttpServer(cfg httpservercfg.Config, authSvc *authentication.Service, userSvc *user.Service, userValidator *uservalidator.Validator) *HttpServer {

	return &HttpServer{serverConfig: cfg, userHandler: userhandler.NewHandler(userSvc, authSvc, userValidator)}
}

func (hs *HttpServer) Serve() {

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/health-check", hs.HealthCheckHandler)

	hs.userHandler.SetRoute(e)

	serverAddress := fmt.Sprintf("%s:%d", hs.serverConfig.ServerCfg.Host, hs.serverConfig.ServerCfg.Port)
	if err := e.Start(serverAddress); err != nil && errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server", "error", err)
	}
}
