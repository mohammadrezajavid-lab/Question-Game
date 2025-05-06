package config

import (
	"golang.project/go-fundamentals/gameapp/repository/mysql"
	"golang.project/go-fundamentals/gameapp/service/authorize"
	"golang.project/go-fundamentals/gameapp/service/user"
	"time"
)

const (
	JWTSignKey            = `jwt_secret_key`
	AccessExpirationTime  = time.Hour * 24
	RefreshExpirationTime = time.Hour * 24 * 7
	AccessSubject         = "at"
	RefreshSubject        = "ar"
	Host                  = "127.0.0.1"
	Port                  = 8080
	DataBaseUserName      = "game_app"
	DataBasePassword      = "game_app_pass"
	DataBaseName          = "game_app_db"
	DataBaseHost          = "127.0.0.1"
	DataBasePort          = 3308
	DataBaseParseTime     = true
)

type SetUpConfig struct {
	Config      Config
	UserService *user.Service
	AuthService *authorize.Service
}

func NewSetUpConfig() SetUpConfig {
	cfg := setUpConfig()
	userSvc, authSvc := setUpSVC(cfg)

	return SetUpConfig{
		Config:      cfg,
		UserService: userSvc,
		AuthService: authSvc,
	}
}

func setUpConfig() Config {

	return NewConfig(
		NewHttpServerCfg(Host, Port),
		authorize.NewConfig(
			[]byte(JWTSignKey),
			AccessExpirationTime,
			RefreshExpirationTime,
			AccessSubject,
			RefreshSubject,
		),
		mysql.NewConfig(
			DataBaseUserName,
			DataBasePassword,
			DataBaseName,
			DataBaseHost,
			DataBaseParseTime,
			DataBasePort,
		),
	)
}

func setUpSVC(cfg Config) (*user.Service, *authorize.Service) {

	authSvc := authorize.NewService(cfg.AuthCfg)
	userSvc := user.NewService(
		mysql.NewDB(
			mysql.NewConfig(
				DataBaseUserName,
				DataBasePassword,
				DataBaseName,
				DataBaseHost,
				DataBaseParseTime,
				DataBasePort,
			),
		),
		authSvc,
	)

	return userSvc, authSvc
}
