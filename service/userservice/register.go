package userservice

import (
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/param/userparam"
	"golang.project/go-fundamentals/gameapp/pkg/hash"
	"golang.project/go-fundamentals/gameapp/pkg/richerror"
)

func (s *Service) Register(req *userparam.RegisterRequest) (*userparam.RegisterResponse, error) {

	const operation = "service.user.Register"

	// TODO - We should verify Phone Number by Verification SMS Code

	hashedPassword, hashErr := hash.Hash(req.Password)
	if hashErr != nil {

		return nil, richerror.NewRichError(operation).
			WithError(hashErr).
			WithMessage("unexpected error: problem for hashing password").
			WithKind(richerror.KindUnexpected)
	}

	user := entity.NewUser(req.Name, req.PhoneNumber, hashedPassword)
	newUser, cErr := s.userRepository.RegisterUser(user)
	if cErr != nil {

		return nil, richerror.NewRichError(operation).
			WithError(cErr).
			WithMessage("unexpected error").
			WithKind(richerror.KindUnexpected).
			WithMeta(map[string]interface{}{"request": req})
	}

	return userparam.NewRegisterResponse(newUser), nil
}
