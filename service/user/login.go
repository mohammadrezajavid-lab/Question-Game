package user

import (
	"golang.project/go-fundamentals/gameapp/dto"
	"golang.project/go-fundamentals/gameapp/pkg/hash"
	"golang.project/go-fundamentals/gameapp/pkg/richerror"
)

func (s *Service) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {

	const operation = "service.user.Login"
	user, gErr := s.userRepository.GetUserByPhoneNumber(req.PhoneNumber)
	if gErr != nil {

		return nil, richerror.NewRichError(operation).WithError(gErr)
	}

	if !hash.CheckHash(req.Password, user.HashedPassword) {

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

	return dto.NewLoginResponse(dto.NewUserInfo(user.ID, user.Name), dto.NewTokens(accessToken, refreshToken)), nil
}
