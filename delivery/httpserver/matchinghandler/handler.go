package matchinghandler

import (
	"golang.project/go-fundamentals/gameapp/adapter/presenceclient"
	"golang.project/go-fundamentals/gameapp/service/authenticationservice"
	"golang.project/go-fundamentals/gameapp/service/authorizationservice"
	"golang.project/go-fundamentals/gameapp/service/matchingservice"
	"golang.project/go-fundamentals/gameapp/validator/matchingvalidator"
)

type MatchingHandler struct {
	authService          authenticationservice.Service
	authorizationService authorizationservice.Service
	matchingService      matchingservice.Service
	matchingValidator    matchingvalidator.Validator
	presenceClient       presenceclient.Client
}

func NewHandler(
	authService authenticationservice.Service,
	authorizationService authorizationservice.Service,
	matchingService matchingservice.Service,
	matchingValidator matchingvalidator.Validator,
	presenceClient presenceclient.Client,
) MatchingHandler {

	return MatchingHandler{
		authService:          authService,
		authorizationService: authorizationService,
		matchingService:      matchingService,
		matchingValidator:    matchingValidator,
		presenceClient:       presenceClient,
	}
}
