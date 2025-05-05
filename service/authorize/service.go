package authorize

import (
	"github.com/golang-jwt/jwt/v4"
	"gocasts.ir/go-fundamentals/gameapp/entity"
	"strings"
	"time"
)

type Service struct {
	signKey               []byte
	accessExpirationTime  time.Duration
	refreshExpirationTime time.Duration
	accessSubject         string
	refreshSubject        string
}

func NewService(
	signKey []byte,
	accessExpirationTime, refreshExpirationTime time.Duration,
	accessSubject, refreshSubject string,
) *Service {

	return &Service{
		signKey:               signKey,
		accessExpirationTime:  accessExpirationTime,
		refreshExpirationTime: refreshExpirationTime,
		accessSubject:         accessSubject,
		refreshSubject:        refreshSubject,
	}
}

func (s *Service) CreateAccessToken(user *entity.User) (string, error) {

	return s.createAccessToken(user.ID)
}

func (s *Service) CreateRefreshToken(user *entity.User) (string, error) {

	return s.createRefreshToken(user.ID)
}

func (s *Service) createAccessToken(userId uint) (string, error) {

	// create a new jwt and set signer SHA 256 in jwt Header
	t := jwt.New(jwt.GetSigningMethod(jwt.SigningMethodHS256.Alg()))

	// set our claims
	t.Claims = NewClaims(
		jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessExpirationTime))},
		s.accessSubject,
		userId,
	)
	// create token string
	return t.SignedString(s.signKey)
}

func (s *Service) ParseJWT(tokenString string) (*Claims, error) {

	// https://pkg.go.dev/github.com/golang-jwt/jwt/v5#example-Parse-Hmac

	tokenString = strings.Replace(tokenString, `Bearer `, "", 1)

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {

		return s.signKey, nil
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

func (s *Service) createRefreshToken(userId uint) (string, error) {

	// create a new jwt and set signer SHA 256 in jwt Header
	t := jwt.New(jwt.GetSigningMethod(jwt.SigningMethodHS256.Alg()))

	// set our claims
	t.Claims = NewClaims(
		jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.refreshExpirationTime))},
		s.refreshSubject,
		userId,
	)

	return t.SignedString(s.signKey)

}
