package userservice

import (
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/metrics"
	"golang.project/go-fundamentals/gameapp/param/userparam"
	"golang.project/go-fundamentals/gameapp/pkg/errormessage"
	"golang.project/go-fundamentals/gameapp/pkg/hash"
	"golang.project/go-fundamentals/gameapp/pkg/richerror"
)

func (s *Service) Register(req *userparam.RegisterRequest) (*userparam.RegisterResponse, error) {

	const operation = "service.user.Register"

	// TODO - We should verify Phone Number by Verification SMS Code

	hashedPassword, hashErr := hash.Hash(req.Password)
	if hashErr != nil {
		metrics.FailedRegisterUserCounter.Inc()
		logger.Warn(hashErr, "hashing password warning")
		logger.Warn(hashErr, "failed register user")

		return nil, richerror.NewRichError(operation).
			WithError(hashErr).
			WithMessage(errormessage.ErrorMsgUnexpected).
			WithKind(richerror.KindUnexpected)
	}

	user := entity.NewUser(req.Name, req.PhoneNumber, hashedPassword)
	newUser, cErr := s.userRepository.RegisterUser(user)
	if cErr != nil {
		metrics.FailedRegisterUserCounter.Inc()
		logger.Warn(cErr, "failed register user")

		return nil, richerror.NewRichError(operation).
			WithError(cErr).
			WithMessage(errormessage.ErrorMsgUnexpected).
			WithKind(richerror.KindUnexpected).
			WithMeta(map[string]interface{}{"request": req})
	}

	return userparam.NewRegisterResponse(newUser), nil
}
