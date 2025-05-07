package authorize

import "github.com/golang-jwt/jwt/v4"

type Claims struct {
	jwt.RegisteredClaims
	Subject string // It always takes these two values. (at or rt) ----> at= accessToken, rt= refreshToken
	UserId  uint   `json:"user_id"`
}

func NewClaims(claims jwt.RegisteredClaims, subject string, userId uint) *Claims {
	return &Claims{
		RegisteredClaims: claims,
		Subject:          subject,
		UserId:           userId,
	}
}
