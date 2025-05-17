package authentication

import (
	"github.com/golang-jwt/jwt/v4"
	"golang.project/go-fundamentals/gameapp/config/httpservercfg/constant"
	"golang.project/go-fundamentals/gameapp/entity"
	"strings"
	"time"
)

// Config The following structure for the auth service config
type Config struct {
	SignKey               string        `mapstructure:"sign_key"`
	AccessExpirationTime  time.Duration `mapstructure:"access_expiration_time"`
	RefreshExpirationTime time.Duration `mapstructure:"refresh_expiration_time"`
	AccessSubject         string        `mapstructure:"access_subject"`
	RefreshSubject        string        `mapstructure:"refresh_subject"`
	//SignMethod            jwt.SigningMethod `mapstructure:"sign_method"`
}

func NewConfig(
	signKey string,
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
	Config Config
}

func NewService(authConfig Config) *Service {

	return &Service{Config: authConfig}
}

func (s *Service) CreateAccessToken(user *entity.User) (string, error) {

	return s.createAccessToken(user.Id)
}

func (s *Service) CreateRefreshToken(user *entity.User) (string, error) {

	return s.createRefreshToken(user.Id)
}

func (s *Service) ParseJWT(tokenString string) (*Claims, error) {

	// https://pkg.go.dev/github.com/golang-jwt/jwt/v5#example-Parse-Hmac

	tokenString = strings.Replace(tokenString, `Bearer `, "", 1)

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {

		return []byte(s.Config.SignKey), nil
	}, jwt.WithValidMethods([]string{constant.DefaultSignMethod}))
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
	t := jwt.New(jwt.GetSigningMethod(constant.DefaultSignMethod))

	// set our claims
	t.Claims = NewClaims(
		jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.Config.AccessExpirationTime))},
		s.Config.AccessSubject,
		userId,
	)

	return t.SignedString([]byte(s.Config.SignKey))
}

func (s *Service) createRefreshToken(userId uint) (string, error) {

	// create a new jwt and set signer SHA 256 in jwt Header
	t := jwt.New(jwt.GetSigningMethod(constant.DefaultSignMethod))

	// set our claims
	t.Claims = NewClaims(
		jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.Config.RefreshExpirationTime))},
		s.Config.RefreshSubject,
		userId,
	)

	return t.SignedString([]byte(s.Config.SignKey))

}
