package authparam

import "golang.project/go-fundamentals/gameapp/entity"

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type ClaimRefreshTokenRequest struct {
	UserId   uint
	UserRole entity.Role
}

type TokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
