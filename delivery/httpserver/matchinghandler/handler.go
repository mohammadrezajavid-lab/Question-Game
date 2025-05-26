package matchinghandler

import (
	"golang.project/go-fundamentals/gameapp/service/authenticationservice"
	"golang.project/go-fundamentals/gameapp/service/authorizationservice"
	"golang.project/go-fundamentals/gameapp/service/matchingservice"
	"golang.project/go-fundamentals/gameapp/service/presenceservice"
	"golang.project/go-fundamentals/gameapp/validator/matchingvalidator"
)

type MatchingHandler struct {
	authService          *authenticationservice.Service
	authorizationService *authorizationservice.Service
	matchingService      *matchingservice.Service
	matchingValidator    *matchingvalidator.Validator
	presenceService      *presenceservice.Service
}

func NewHandler(
	authService *authenticationservice.Service,
	authorizationService *authorizationservice.Service,
	matchingService *matchingservice.Service,
	matchingValidator *matchingvalidator.Validator,
	presenceService *presenceservice.Service,
) MatchingHandler {

	return MatchingHandler{
		authService:          authService,
		authorizationService: authorizationService,
		matchingService:      matchingService,
		matchingValidator:    matchingValidator,
		presenceService:      presenceService,
	}
}
