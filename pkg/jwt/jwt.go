package jwt

import (
	"github.com/golang-jwt/jwt/v4"
	"golang.project/go-fundamentals/gameapp/entity"
	"strings"
	"time"
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

type Config struct {
	SignKey               string        `mapstructure:"sign_key"`
	AccessExpirationTime  time.Duration `mapstructure:"access_expiration_time"`
	RefreshExpirationTime time.Duration `mapstructure:"refresh_expiration_time"`
	AccessSubject         string        `mapstructure:"access_subject"`
	RefreshSubject        string        `mapstructure:"refresh_subject"`
	SignMethod            string        `mapstructure:"sign_method"`
}
type JWT struct {
	Config Config
}

func NewJWT(config Config) *JWT {
	return &JWT{Config: config}
}

func (j *JWT) ParseJWT(tokenString string) (*Claims, error) {
	tokenString = strings.Replace(tokenString, `Bearer `, "", 1)

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.Config.SignKey), nil
	}, jwt.WithValidMethods([]string{j.Config.SignMethod}))
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok {
		return claims, nil
	} else {
		return nil, err
	}
}

func (j *JWT) CreateAccessToken(userId uint, role entity.Role) (string, error) {

	t := jwt.New(jwt.GetSigningMethod(j.Config.SignMethod))

	t.Claims = NewClaims(
		jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.Config.AccessExpirationTime))},
		j.Config.AccessSubject,
		userId,
		role,
	)

	return t.SignedString([]byte(j.Config.SignKey))
}

func (j *JWT) CreateRefreshToken(userId uint, role entity.Role) (string, error) {

	t := jwt.New(jwt.GetSigningMethod(j.Config.SignMethod))

	t.Claims = NewClaims(
		jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.Config.RefreshExpirationTime))},
		j.Config.RefreshSubject,
		userId,
		role,
	)

	return t.SignedString([]byte(j.Config.SignKey))
}
