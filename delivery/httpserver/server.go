package httpserver

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.project/go-fundamentals/gameapp/config/httpservercfg"
	"golang.project/go-fundamentals/gameapp/delivery/httpserver/backofficeuserhandler"
	"golang.project/go-fundamentals/gameapp/delivery/httpserver/matchinghandler"
	"golang.project/go-fundamentals/gameapp/delivery/httpserver/userhandler"
	"golang.project/go-fundamentals/gameapp/service/authenticationservice"
	"golang.project/go-fundamentals/gameapp/service/authorizationservice"
	"golang.project/go-fundamentals/gameapp/service/backofficeuserservice"
	"golang.project/go-fundamentals/gameapp/service/matchingservice"
	"golang.project/go-fundamentals/gameapp/service/userservice"
	"golang.project/go-fundamentals/gameapp/validator/matchingvalidator"
	"golang.project/go-fundamentals/gameapp/validator/uservalidator"
	"log/slog"
	"net/http"
)

type HttpServer struct {
	config httpservercfg.Config

	userHandler           userhandler.UserHandler
	backOfficeUserHandler backofficeuserhandler.BackOfficeUserHandler
	matchingHandler       matchinghandler.MatchingHandler
}

func NewHttpServer(
	cfg httpservercfg.Config,
	authSvc *authenticationservice.Service,
	userSvc *userservice.Service,
	backOfficeUserSvc *backofficeuserservice.Service,
	authorizationSvc *authorizationservice.Service,
	userValidator *uservalidator.Validator,
	matchingSvc *matchingservice.Service,
	matchingValidator *matchingvalidator.Validator,
) *HttpServer {

	return &HttpServer{
		config: cfg,

		userHandler:           userhandler.NewHandler(userSvc, authSvc, authorizationSvc, userValidator),
		backOfficeUserHandler: backofficeuserhandler.NewHandler(backOfficeUserSvc, authSvc, authorizationSvc, userValidator),
		matchingHandler:       matchinghandler.NewHandler(authSvc, authorizationSvc, matchingSvc, matchingValidator),
	}
}

func (hs *HttpServer) Serve() {

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/health-check", hs.HealthCheckHandler)

	hs.userHandler.SetRoute(e)
	hs.backOfficeUserHandler.SetRoute(e)
	hs.matchingHandler.SetRoute(e)

	serverAddress := fmt.Sprintf("%s:%d", hs.config.ServerCfg.Host, hs.config.ServerCfg.Port)
	if err := e.Start(serverAddress); err != nil && errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server", "error", err)
	}
}
