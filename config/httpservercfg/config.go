package httpservercfg

import (
	"database/sql"
	"fmt"
	"golang.project/go-fundamentals/gameapp/repository/migrator"
	"golang.project/go-fundamentals/gameapp/repository/mysql"
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
	//SignMethod            = jwt.SigningMethodHS256.Alg()
)

type Config struct {
	Host           string
	Port           int
	DataBaseConfig mysql.Config
}

func NewConfig(host string, port int) Config {

	dbConf := mysql.NewConfig(
		DataBaseUserName,
		DataBasePassword,
		DataBaseName,
		DataBaseHost,
		DataBaseParseTime,
		DataBasePort,
	)

	return Config{
		Host:           host,
		Port:           port,
		DataBaseConfig: dbConf,
	}
}

func (c Config) SetUpConfig(migrationCommand string) {

	if migrationCommand != "up" && migrationCommand != "down" && migrationCommand != "skip" && migrationCommand != "status" {
		panic(fmt.Sprintf("invalid migration-command: %s", migrationCommand))
	}

	dbConnection := mysql.NewDB(c.DataBaseConfig).MysqlConnection
	c.migrate(dbConnection, MigrateDialect, migrationCommand)

}

func (c Config) migrate(dbConnection *sql.DB, dialect string, migrationCommand string) {

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
