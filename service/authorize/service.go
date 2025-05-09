package authorize

import (
	"github.com/golang-jwt/jwt/v4"
	"golang.project/go-fundamentals/gameapp/entity"
	"strings"
	"time"
)

// Config The following structure for the auth service config
type Config struct {
	SignKey               []byte
	AccessExpirationTime  time.Duration
	RefreshExpirationTime time.Duration
	AccessSubject         string
	RefreshSubject        string
}

func NewConfig(
	signKey []byte,
	accessExpirationTime, refreshExpirationTime time.Duration,
	accessSubject, refreshSubject string,
) Config {

	return Config{
		SignKey:               signKey,
		AccessExpirationTime:  accessExpirationTime,
		RefreshExpirationTime: refreshExpirationTime,
		AccessSubject:         accessSubject,
		RefreshSubject:        refreshSubject,
	}
}

type Service struct {
	config Config
}

func NewService(config Config) *Service {

	return &Service{config: config}
}

func (s *Service) CreateAccessToken(user *entity.User) (string, error) {

	return s.createAccessToken(user.ID)
}

func (s *Service) CreateRefreshToken(user *entity.User) (string, error) {

	return s.createRefreshToken(user.ID)
}

func (s *Service) ParseJWT(tokenString string) (*Claims, error) {

	// https://pkg.go.dev/github.com/golang-jwt/jwt/v5#example-Parse-Hmac

	tokenString = strings.Replace(tokenString, `Bearer `, "", 1)

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {

		return s.config.SignKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {

		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok {

		return claims, nil
	} else {

		return nil, err
	}
}

func (s *Service) createAccessToken(userId uint) (string, error) {

	// create a new jwt and set signer SHA 256 in jwt Header
	t := jwt.New(jwt.GetSigningMethod(jwt.SigningMethodHS256.Alg()))

	// set our claims
	t.Claims = NewClaims(
		jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.AccessExpirationTime))},
		s.config.AccessSubject,
		userId,
	)
	// create token string
	return t.SignedString(s.config.SignKey)
}

func (s *Service) createRefreshToken(userId uint) (string, error) {

	// create a new jwt and set signer SHA 256 in jwt Header
	t := jwt.New(jwt.GetSigningMethod(jwt.SigningMethodHS256.Alg()))

	// set our claims
	t.Claims = NewClaims(
		jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.RefreshExpirationTime))},
		s.config.RefreshSubject,
		userId,
	)

	return t.SignedString(s.config.SignKey)

}
