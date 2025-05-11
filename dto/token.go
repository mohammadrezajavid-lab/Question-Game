package dto

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewTokens(accessToken, refreshToken string) Tokens {
	return Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
