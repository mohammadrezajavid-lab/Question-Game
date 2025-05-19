package userservice

import (
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/param"
	"golang.project/go-fundamentals/gameapp/pkg/hash"
	"golang.project/go-fundamentals/gameapp/pkg/richerror"
)

func (s *Service) Register(req *param.RegisterRequest) (*param.RegisterResponse, error) {

	const operation = "service.user.Register"

	// TODO - We should verify Phone Number by Verification SMS Code

	hashedPassword, hashErr := hash.Hash(req.Password)
	if hashErr != nil {

		return nil, richerror.NewRichError(operation).
			WithError(hashErr).
			WithMessage("unexpected error: problem for hashing password").
			WithKind(richerror.KindUnexpected)
	}

	// create new user in storage
	user := entity.NewUser(req.Name, req.PhoneNumber, hashedPassword)
	//TODO - user.Role = entity.AdminRole // for change default user Role
	newUser, cErr := s.userRepository.RegisterUser(user)
	if cErr != nil {

		return nil, richerror.NewRichError(operation).
			WithError(cErr).
			WithMessage("unexpected error").
			WithKind(richerror.KindUnexpected).
			WithMeta(map[string]interface{}{"request": req})
	}

	// return created user
	return param.NewRegisterResponse(newUser), nil
}
