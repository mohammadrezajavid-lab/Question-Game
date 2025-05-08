package user

import (
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/pkg/password"
	"golang.project/go-fundamentals/gameapp/pkg/phonenumber"
	"golang.project/go-fundamentals/gameapp/pkg/richerror"
	"golang.project/go-fundamentals/gameapp/service/contract"
)

type Service struct {
	userRepository contract.UserRepository
	authService    contract.AuthorizeGenerator
}

func NewService(userRepository contract.UserRepository, authorizeService contract.AuthorizeGenerator) *Service {
	return &Service{
		userRepository: userRepository,
		authService:    authorizeService,
	}
}

func (s *Service) Register(req *RegisterRequest) (*RegisterResponse, error) {

	const operation = "service.user.Register"

	// TODO - We should verify Phone Number by Verification SMS Code

	// validate phone number
	if !phonenumber.IsPhoneNumberValid(req.PhoneNumber) {

		return nil, richerror.NewRichError(operation).
			WithMessage("phone number is invalid").
			WithKind(richerror.KindInvalid).
			WithMeta(map[string]interface{}{"phone_number": req.PhoneNumber})
	}

	// check uniqueness of phone number
	if isUniq, err := s.userRepository.IsPhoneNumberUniq(req.PhoneNumber); !isUniq || err != nil {
		if !isUniq {

			return nil, richerror.NewRichError(operation).
				WithError(err).
				WithMessage("phone number is not uniq").
				WithKind(richerror.KindInvalid).
				WithMeta(map[string]interface{}{"phone_number": req.PhoneNumber})
		}

		return nil, richerror.NewRichError(operation).
			WithError(err).
			WithMessage("unexpected error").
			WithKind(richerror.KindUnexpected)
	}

	// validate name
	if !phonenumber.IsNameValid(req.Name) {

		return nil, richerror.NewRichError(operation).
			WithMessage("name length should be greater than 3").
			WithKind(richerror.KindInvalid).
			WithMeta(map[string]interface{}{"name": req.Name})
	}

	// validate password
	// TODO - It is better to use Regex for password.
	if len(req.Password) < 8 {

		return nil, richerror.NewRichError(operation).
			WithMessage("password length should be greater than 8").
			WithKind(richerror.KindInvalid).
			WithMeta(map[string]interface{}{"password": req.Password})
	}

	hashedPassword, hashErr := password.HashPassword(req.Password)
	if hashErr != nil {

		return nil, richerror.NewRichError(operation).
			WithError(hashErr).
			WithMessage("unexpected error: problem for hashing password").
			WithKind(richerror.KindUnexpected)
	}

	// create new user in storage
	newUser, cErr := s.userRepository.RegisterUser(entity.NewUser(req.Name, req.PhoneNumber, hashedPassword))
	if cErr != nil {

		return nil, richerror.NewRichError(operation).
			WithError(cErr).
			WithMessage("unexpected error").
			WithKind(richerror.KindUnexpected).
			WithMeta(map[string]interface{}{"request": req})
	}

	// return created user
	return NewRegisterResponse(newUser), nil
}

func (s *Service) Login(req *LoginRequest) (*LoginResponse, error) {

	const operation = "service.user.Login"
	user, exist, gErr := s.userRepository.GetUserByPhoneNumber(req.PhoneNumber)
	if gErr != nil {

		if !exist {

			// error: record not found
			return nil, richerror.NewRichError(operation).
				WithError(gErr).
				WithMessage("phoneNumber or password incorrect").
				WithKind(richerror.KindNotFound)
		}

		return nil, richerror.NewRichError(operation).
			WithError(gErr).
			WithMessage("unexpected error").
			WithKind(richerror.KindUnexpected)
	}

	if !password.CheckPasswordHash(req.Password, user.HashedPassword) {

		return nil, richerror.NewRichError(operation).
			WithMessage("phoneNumber or password incorrect").
			WithKind(richerror.KindNotFound).
			WithMeta(map[string]interface{}{"request": req})
	}

	// If the user exists create accessToken and refreshToken

	accessToken, aErr := s.authService.CreateAccessToken(user)
	if aErr != nil {

		return nil, richerror.NewRichError(operation).
			WithError(aErr).
			WithMessage("unexpected error").
			WithKind(richerror.KindUnexpected)
	}

	refreshToken, rErr := s.authService.CreateRefreshToken(user)
	if rErr != nil {

		return nil, richerror.NewRichError(operation).
			WithError(rErr).
			WithMessage("unexpected error").
			WithKind(richerror.KindUnexpected)
	}

	return NewLoginResponse(NewUserInfo(user.ID, user.Name), NewTokens(accessToken, refreshToken)), nil
}

// All Request Inputs for Interaction/Service Should be Sanitized.

func (s *Service) Profile(req *ProfileRequest) (*ProfileResponse, error) {

	const operation = "service.user.Profile"
	user, err := s.userRepository.GetUserById(req.UserId)
	if err != nil {

		return nil, richerror.NewRichError(operation).WithError(err)
	}

	return NewProfileResponse(user.Name), nil
}
