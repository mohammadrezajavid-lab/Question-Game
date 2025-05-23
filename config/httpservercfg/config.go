package httpservercfg

import (
	"database/sql"
	"fmt"
	"github.com/spf13/viper"
	"golang.project/go-fundamentals/gameapp/adapter/redis"
	"golang.project/go-fundamentals/gameapp/config/httpservercfg/constant"
	"golang.project/go-fundamentals/gameapp/repository/migrator"
	"golang.project/go-fundamentals/gameapp/repository/mysql"
	"golang.project/go-fundamentals/gameapp/service/authenticationservice"
	"golang.project/go-fundamentals/gameapp/service/matchingservice"
	"log"
	"strings"
)

type HttpServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type Config struct {
	ServerCfg   HttpServerConfig             `mapstructure:"httpserver_cfg"`
	DataBaseCfg mysql.Config                 `mapstructure:"database_cfg"`
	AuthCfg     authenticationservice.Config `mapstructure:"auth_cfg"`
	MatchingCfg matchingservice.Config       `mapstructure:"matching_cfg"`
	RedisCfg    redis.Config                 `mapstructure:"redis_cfg"`
}

func NewConfig(host string, port int) Config {

	appCfg := loadConfig(host, port)

	return Config{
		ServerCfg:   appCfg.ServerCfg,
		DataBaseCfg: appCfg.DataBaseCfg,
		AuthCfg:     appCfg.AuthCfg,
		MatchingCfg: appCfg.MatchingCfg,
		RedisCfg:    appCfg.RedisCfg,
	}
}

// 1. read config file
// 2. env variable
// 3. use default env
func loadConfig(host string, port int) Config {

	setDefaultENV()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// read config file
	viper.SetConfigName(constant.DefaultConfigFileName)
	viper.SetConfigType(constant.DefaultConfigFileType)
	viper.AddConfigPath(constant.DefaultConfigFilePath)

	var appConfig Config
	if err := viper.ReadInConfig(); err != nil {

		log.Println("⚠️ config file not found, using environment variables or default values.")

		// get config from env variable
		if uErr := viper.Sub("httpserver_cfg").Unmarshal(&appConfig.ServerCfg); uErr != nil {
			log.Fatalf("can't unmarshal httpserver config: %v", uErr)
		}
		if uErr := viper.Sub("database_cfg").Unmarshal(&appConfig.DataBaseCfg); uErr != nil {
			log.Fatalf("can't unmarshal database config: %v", uErr)
		}
		if uErr := viper.Sub("auth_cfg").Unmarshal(&appConfig.AuthCfg); uErr != nil {
			log.Fatalf("can't unmarshal auth config: %v", uErr)
		}
		if uErr := viper.Sub("matching_cfg").Unmarshal(&appConfig.MatchingCfg); uErr != nil {
			log.Fatalf("can't unmarshal matching config: %v", uErr)
		}
		if uErr := viper.Sub("redis_cfg").Unmarshal(&appConfig.RedisCfg); uErr != nil {
			log.Fatalf("can't unmarshal redis config: %v", uErr)
		}

	} else {

		if uErr := viper.Unmarshal(&appConfig); uErr != nil {

			panic(fmt.Errorf("can't Unmarshal config file into struct Config, %w", uErr))
		}
	}

	if host != "" {
		appConfig.ServerCfg.Host = host
	}
	if port != 0 {
		appConfig.ServerCfg.Port = port
	}

	return appConfig
}

func setDefaultENV() {

	// Default HTTP Server config ENV
	viper.SetDefault("httpserver_cfg.host", constant.DefaultHTTPServerHost)
	viper.SetDefault("httpserver_cfg.port", constant.DefaultHTTPServerPort)

	// Default DataBase config ENV
	viper.SetDefault("database_cfg.database_user_name", constant.DefaultDataBaseUserName)
	viper.SetDefault("database_cfg.database_password", constant.DefaultDataBasePassword)
	viper.SetDefault("database_cfg.database_name", constant.DefaultDataBaseName)
	viper.SetDefault("database_cfg.database_host", constant.DefaultDataBaseHost)
	viper.SetDefault("database_cfg.database_parse_time", constant.DefaultDataBaseParseTime)
	viper.SetDefault("database_cfg.database_port", constant.DefaultDataBasePort)

	// Default Auth config ENV
	viper.SetDefault("auth_cfg.sign_key", constant.DefaultJWTSignKey)
	viper.SetDefault("auth_cfg.access_expiration_time", constant.DefaultAccessExpirationTime)
	viper.SetDefault("auth_cfg.refresh_expiration_time", constant.DefaultRefreshExpirationTime)
	viper.SetDefault("auth_cfg.access_subject", constant.DefaultAccessSubject)
	viper.SetDefault("auth_cfg.refresh_subject", constant.DefaultRefreshSubject)

	// Default Matching config ENV
	viper.SetDefault("waiting_time_out", constant.DefaultWaitingTimeOut)
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
