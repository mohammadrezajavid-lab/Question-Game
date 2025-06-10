package userservice

import (
	"golang.project/go-fundamentals/gameapp/metrics"
	"golang.project/go-fundamentals/gameapp/param/userparam"
	"golang.project/go-fundamentals/gameapp/pkg/errormessage"
	"golang.project/go-fundamentals/gameapp/pkg/hash"
	"golang.project/go-fundamentals/gameapp/pkg/richerror"
)

func (s *Service) Login(req *userparam.LoginRequest) (*userparam.LoginResponse, error) {

	const operation = "service.user.Login"
	user, gErr := s.userRepository.GetUserByPhoneNumber(req.PhoneNumber)
	if gErr != nil {
		metrics.FailedLoginCounter.Inc()

		return nil, richerror.NewRichError(operation).WithError(gErr)
	}

	if !hash.CheckHash(req.Password, user.HashedPassword) {
		metrics.FailedLoginCounter.Inc()
		metrics.FailedLoginIncorrectPhoneNumberOrPasswordCounter.Inc()

		return nil, richerror.NewRichError(operation).
			WithMessage(errormessage.ErrorMsgIncorrectPhoneNumberPassword).
			WithKind(richerror.KindNotFound).
			WithMeta(map[string]interface{}{"request": req})
	}

	// If the user exists create accessToken and refreshToken

	accessToken, aErr := s.authService.CreateAccessToken(user)
	if aErr != nil {
		metrics.FailedCreateAccessTokenCounter.Inc()

		return nil, richerror.NewRichError(operation).
			WithError(aErr).
			WithMessage(errormessage.ErrorMsgUnexpected).
			WithKind(richerror.KindUnexpected)
	}

	refreshToken, rErr := s.authService.CreateRefreshToken(user)
	if rErr != nil {
		metrics.FailedCreateRefreshTokenCounter.Inc()

		return nil, richerror.NewRichError(operation).
			WithError(rErr).
			WithMessage(errormessage.ErrorMsgUnexpected).
			WithKind(richerror.KindUnexpected)
	}

	return userparam.NewLoginResponse(userparam.NewUserInfo(user.Id, user.Name), userparam.NewTokens(accessToken, refreshToken)), nil
}
