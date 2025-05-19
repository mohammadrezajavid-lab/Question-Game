package authenticationservice

import (
	"github.com/golang-jwt/jwt/v4"
	"golang.project/go-fundamentals/gameapp/entity"
)

type Claims struct {
	jwt.RegisteredClaims
	Subject string      // It always takes these two values. (at or rt) ----> at= accessToken, rt= refreshToken
	UserId  uint        `json:"user_id"`
	Role    entity.Role `json:"role"`
}

func NewClaims(claims jwt.RegisteredClaims, subject string, userId uint, role entity.Role) *Claims {
	return &Claims{
		RegisteredClaims: claims,
		Subject:          subject,
		UserId:           userId,
		Role:             role,
	}
}
