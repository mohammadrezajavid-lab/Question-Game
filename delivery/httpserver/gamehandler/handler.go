package gamehandler

import (
	"golang.project/go-fundamentals/gameapp/adapter/presenceclient"
	"golang.project/go-fundamentals/gameapp/service/authenticationservice"
	"golang.project/go-fundamentals/gameapp/service/authorizationservice"
	"golang.project/go-fundamentals/gameapp/service/gameservice"
	"golang.project/go-fundamentals/gameapp/service/matchingservice"
	"golang.project/go-fundamentals/gameapp/validator/matchingvalidator"
)

type GameHandler struct {
	authService          authenticationservice.Service
	authorizationService authorizationservice.Service
	matchingService      matchingservice.Service
	matchingValidator    matchingvalidator.Validator
	presenceClient       presenceclient.Client
	gameService          gameservice.Service
}

func NewHandler(
	authService authenticationservice.Service,
	authorizationService authorizationservice.Service,
	matchingService matchingservice.Service,
	matchingValidator matchingvalidator.Validator,
	presenceClient presenceclient.Client,
	gameService gameservice.Service,
) GameHandler {

	return GameHandler{
		authService:          authService,
		authorizationService: authorizationService,
		matchingService:      matchingService,
		matchingValidator:    matchingValidator,
		presenceClient:       presenceClient,
		gameService:          gameService,
	}
}
