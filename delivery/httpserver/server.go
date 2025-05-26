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
	"golang.project/go-fundamentals/gameapp/service/presenceservice"
	"golang.project/go-fundamentals/gameapp/service/userservice"
	"golang.project/go-fundamentals/gameapp/validator/matchingvalidator"
	"golang.project/go-fundamentals/gameapp/validator/uservalidator"
	"log"
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
	authSvc *authenticationservice.Service,
	userSvc *userservice.Service,
	backOfficeUserSvc *backofficeuserservice.Service,
	authorizationSvc *authorizationservice.Service,
	userValidator *uservalidator.Validator,
	matchingSvc *matchingservice.Service,
	matchingValidator *matchingvalidator.Validator,
	presenceSvc *presenceservice.Service,
) *Server {

	return &Server{
		config:                cfg,
		userHandler:           userhandler.NewHandler(userSvc, authSvc, authorizationSvc, userValidator, presenceSvc),
		backOfficeUserHandler: backofficeuserhandler.NewHandler(backOfficeUserSvc, authSvc, authorizationSvc, userValidator, presenceSvc),
		matchingHandler:       matchinghandler.NewHandler(authSvc, authorizationSvc, matchingSvc, matchingValidator, presenceSvc),
		router:                echo.New(),
	}
}

func (s *Server) GetRouter() *echo.Echo {
	return s.router
}

func (s *Server) Serve() {

	// Middleware
	s.router.Use(middleware.Logger())
	s.router.Use(middleware.Recover())

	s.router.GET("/health-check", s.HealthCheckHandler)

	s.userHandler.SetRoute(s.router)
	s.backOfficeUserHandler.SetRoute(s.router)
	s.matchingHandler.SetRoute(s.router)

	serverAddress := fmt.Sprintf("%s:%d", s.config.ServerCfg.Host, s.config.ServerCfg.Port)
	if err := s.router.Start(serverAddress); err != nil && errors.Is(err, http.ErrServerClosed) {
		log.Printf("failed to start server, error: %v\n", err)
	}
}
