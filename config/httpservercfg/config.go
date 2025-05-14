package httpservercfg

import (
	"database/sql"
	"fmt"
	"github.com/spf13/viper"
	"golang.project/go-fundamentals/gameapp/config/httpservercfg/constant"
	"golang.project/go-fundamentals/gameapp/repository/migrator"
	"golang.project/go-fundamentals/gameapp/repository/mysql"
	"golang.project/go-fundamentals/gameapp/service/auth"
	"log"
	"strings"
)

type HttpServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type Config struct {
	ServerCfg   HttpServerConfig `mapstructure:"httpserver_cfg"`
	DataBaseCfg mysql.Config     `mapstructure:"database_cfg"`
	AuthCfg     auth.Config      `mapstructure:"auth_cfg"`
}

func NewConfig(host string, port int) Config {

	serverCfg, dataBaseCfg, authCfg := loadConfig(host, port)

	return Config{
		ServerCfg:   serverCfg,
		DataBaseCfg: dataBaseCfg,
		AuthCfg:     authCfg,
	}
}

func loadConfig(host string, port int) (HttpServerConfig, mysql.Config, auth.Config) {

	setDefaultENV()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetConfigName(constant.DefaultConfigFileName)
	viper.SetConfigType(constant.DefaultConfigFileType)
	viper.AddConfigPath(constant.DefaultConfigFilePath)

	// read config file
	if err := viper.ReadInConfig(); err != nil {

		log.Println("⚠️ config file not found, using environment variables or default values.")
	}
	var appConfig Config

	if err := viper.Unmarshal(&appConfig); err != nil {

		panic(fmt.Errorf("can't Unmarshal config file into struct Config, %w", err))
	}

	if host != "" {
		appConfig.ServerCfg.Host = host
	}
	if port != 0 {
		appConfig.ServerCfg.Port = port
	}

	return appConfig.ServerCfg, appConfig.DataBaseCfg, appConfig.AuthCfg
}

func setDefaultENV() {

	// Default HTTP Server ENV
	viper.SetDefault("httpserver_cfg.host", constant.DefaultHTTPServerHost)
	viper.SetDefault("httpserver_cfg.port", constant.DefaultHTTPServerPort)

	// Default DataBase ENV
	viper.SetDefault("database_cfg.database_user_name", constant.DefaultDataBaseUserName)
	viper.SetDefault("database_cfg.database_password", constant.DefaultDataBasePassword)
	viper.SetDefault("database_cfg.database_name", constant.DefaultDataBaseName)
	viper.SetDefault("database_cfg.database_host", constant.DefaultDataBaseHost)
	viper.SetDefault("database_cfg.database_parse_time", constant.DefaultDataBaseParseTime)
	viper.SetDefault("database_cfg.database_port", constant.DefaultDataBasePort)

	// Default Auth ENV
	viper.SetDefault("auth_cfg.sign_key", constant.DefaultJWTSignKey)
	viper.SetDefault("auth_cfg.access_expiration_time", constant.DefaultAccessExpirationTime)
	viper.SetDefault("auth_cfg.refresh_expiration_time", constant.DefaultRefreshExpirationTime)
	viper.SetDefault("auth_cfg.access_subject", constant.DefaultAccessSubject)
	viper.SetDefault("auth_cfg.refresh_subject", constant.DefaultRefreshSubject)
}

func (c *Config) Migrate(migrationCommand string) {

	if migrationCommand != "up" && migrationCommand != "down" && migrationCommand != "skip" && migrationCommand != "status" {
		panic(fmt.Sprintf("invalid migration-command: %s", migrationCommand))
	}

	dbConnection := mysql.NewDB(c.DataBaseCfg).MysqlConnection
	c.migrate(dbConnection, constant.MigrateDialect, migrationCommand)

}

func (c *Config) migrate(dbConnection *sql.DB, dialect string, migrationCommand string) {

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
