package config

import (
	"database/sql"
	"fmt"
	"golang.project/go-fundamentals/gameapp/repository/migrator"
	"golang.project/go-fundamentals/gameapp/repository/mysql"
	"golang.project/go-fundamentals/gameapp/service/auth"
	"golang.project/go-fundamentals/gameapp/service/user"
	"golang.project/go-fundamentals/gameapp/validator/uservalidator"
	"time"
)

const (
	JWTSignKey            = `jwt_secret_key`
	AccessExpirationTime  = time.Hour * 24
	RefreshExpirationTime = time.Hour * 24 * 7
	AccessSubject         = "at"
	RefreshSubject        = "ar"
	DataBaseUserName      = "game_app"
	DataBasePassword      = "game_app_pass"
	DataBaseName          = "game_app_db"
	DataBaseHost          = "127.0.0.1"
	DataBasePort          = 3308
	DataBaseParseTime     = true
	MigrateDialect        = "mysql"
)

type SetUpConfig struct {
	Config Config

	UserService *user.Service
	AuthService *auth.Service

	UserValidator *uservalidator.Validator
}

func NewSetUpConfig(host string, port int, migrationCommand string) SetUpConfig {

	if migrationCommand != "up" && migrationCommand != "down" && migrationCommand != "skip" && migrationCommand != "status" {
		panic(fmt.Sprintf("invalid migration-command: %s", migrationCommand))
	}

	cfg := setUpConfig(host, port)
	repository := mysql.NewDB(cfg.DataBaseCfg)
	userSvc, authSvc := setUpSVC(cfg, repository)

	userValidator := uservalidator.NewValidator(repository)

	setUpMigration(mysql.NewDB(cfg.DataBaseCfg).MysqlConnection, MigrateDialect, migrationCommand)

	return SetUpConfig{
		Config:        cfg,
		UserService:   userSvc,
		AuthService:   authSvc,
		UserValidator: userValidator,
	}
}

func setUpConfig(host string, port int) Config {

	return NewConfig(
		NewHttpServerCfg(host, port),
		auth.NewConfig(
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

func setUpSVC(cfg Config, repository *mysql.DB) (*user.Service, *auth.Service) {

	authSvc := auth.NewService(cfg.AuthCfg)
	userSvc := user.NewService(repository, authSvc)

	return userSvc, authSvc
}

func setUpMigration(dbConnection *sql.DB, dialect string, migrationCommand string) {

	mgt := migrator.NewMigrator(dbConnection, dialect)

	switch migrationCommand {
	case "up":
		mgt.Up()
	case "down":
		mgt.Down()
	case "status":
		mgt.Status()
	default:
		return
	}
}
