package authenticationservice

import (
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/param/authparam"
	"golang.project/go-fundamentals/gameapp/pkg/jwt"
)

type Service struct {
	Jwt *jwt.JWT
}

func NewService(jwt *jwt.JWT) Service {
	return Service{Jwt: jwt}
}

func (s *Service) CreateAccessToken(user *entity.User) (string, error) {

	return s.Jwt.CreateAccessToken(user.Id, user.Role)
}

func (s *Service) CreateRefreshToken(user *entity.User) (string, error) {

	return s.Jwt.CreateRefreshToken(user.Id, user.Role)
}

func (s *Service) CreateTokens(request *authparam.ClaimRefreshTokenRequest) (*authparam.TokensResponse, error) {
	accessToken, caErr := s.Jwt.CreateAccessToken(request.UserId, request.UserRole)
	if caErr != nil {
		return nil, caErr
	}
	refreshToken, crErr := s.Jwt.CreateRefreshToken(request.UserId, request.UserRole)
	if crErr != nil {
		return nil, crErr
	}

	return &authparam.TokensResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
