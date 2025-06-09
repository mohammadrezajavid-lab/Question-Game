package httpserver

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.project/go-fundamentals/gameapp/adapter/presenceclient"
	"golang.project/go-fundamentals/gameapp/config/httpservercfg"
	"golang.project/go-fundamentals/gameapp/delivery/httpserver/backofficeuserhandler"
	"golang.project/go-fundamentals/gameapp/delivery/httpserver/matchinghandler"
	middleware2 "golang.project/go-fundamentals/gameapp/delivery/httpserver/middleware"
	"golang.project/go-fundamentals/gameapp/delivery/httpserver/userhandler"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/pkg/errormessage"
	"golang.project/go-fundamentals/gameapp/service/authenticationservice"
	"golang.project/go-fundamentals/gameapp/service/authorizationservice"
	"golang.project/go-fundamentals/gameapp/service/backofficeuserservice"
	"golang.project/go-fundamentals/gameapp/service/matchingservice"
	"golang.project/go-fundamentals/gameapp/service/userservice"
	"golang.project/go-fundamentals/gameapp/validator/matchingvalidator"
	"golang.project/go-fundamentals/gameapp/validator/uservalidator"
	"net/http"
)

type Server struct {
	config                httpservercfg.Config
	userHandler           userhandler.UserHandler
	backOfficeUserHandler backofficeuserhandler.BackOfficeUserHandler
	matchingHandler       matchinghandler.MatchingHandler
	router                *echo.Echo
}

func New(
	cfg httpservercfg.Config,
	authSvc authenticationservice.Service,
	userSvc userservice.Service,
	backOfficeUserSvc backofficeuserservice.Service,
	authorizationSvc authorizationservice.Service,
	userValidator uservalidator.Validator,
	matchingSvc matchingservice.Service,
	matchingValidator matchingvalidator.Validator,
	presenceClient presenceclient.Client,
) *Server {

	return &Server{
		config:                cfg,
		userHandler:           userhandler.NewHandler(userSvc, authSvc, authorizationSvc, userValidator, presenceClient),
		backOfficeUserHandler: backofficeuserhandler.NewHandler(backOfficeUserSvc, authSvc, authorizationSvc, userValidator, presenceClient),
		matchingHandler:       matchinghandler.NewHandler(authSvc, authorizationSvc, matchingSvc, matchingValidator, presenceClient),
		router:                echo.New(),
	}
}

func (s *Server) GetRouter() *echo.Echo {
	return s.router
}

func (s *Server) Serve() {

	s.router.Use(middleware.RequestID())
	s.router.Use(middleware2.ZapLogger())
	s.router.Use(middleware2.PrometheusMiddleware())
	s.router.Use(middleware.Recover())

	s.router.GET("/health-check", s.HealthCheckHandler)

	s.userHandler.SetRoute(s.router)
	s.backOfficeUserHandler.SetRoute(s.router)
	s.matchingHandler.SetRoute(s.router)

	serverAddress := fmt.Sprintf("%s:%d", s.config.ServerCfg.Host, s.config.ServerCfg.Port)
	if err := s.router.Start(serverAddress); err != nil && errors.Is(err, http.ErrServerClosed) {
		logger.Warn(err, errormessage.ErrorMsgFailedStartServer)
	}
}
