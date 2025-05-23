package constant

import (
	"time"
)

const (
	DefaultConfigFileName = "config"
	DefaultConfigFileType = "yaml"
	DefaultConfigFilePath = "."

	DefaultHTTPServerHost = "127.0.0.1"
	DefaultHTTPServerPort = 8080

	DefaultJWTSignKey            = `jwt_secret_key`
	DefaultSignMethod            = "HS256"
	DefaultAccessExpirationTime  = time.Hour * 24
	DefaultRefreshExpirationTime = time.Hour * 24 * 7
	DefaultAccessSubject         = "at"
	DefaultRefreshSubject        = "rt"

	DefaultDataBaseUserName  = "game_app"
	DefaultDataBasePassword  = "game_app_pass"
	DefaultDataBaseName      = "game_app_db"
	DefaultDataBaseHost      = "127.0.0.1"
	DefaultDataBasePort      = 3308
	DefaultDataBaseParseTime = true

	MigrateDialect           = "mysql"
	AuthMiddlewareContextKey = "claims"

	DefaultWaitingTimeOut = time.Minute * 2
)
