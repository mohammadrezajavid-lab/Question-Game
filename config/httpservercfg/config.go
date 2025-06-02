package httpservercfg

import (
	"database/sql"
	"fmt"
	"github.com/spf13/viper"
	"golang.project/go-fundamentals/gameapp/adapter/presenceclient"
	"golang.project/go-fundamentals/gameapp/adapter/redis"
	"golang.project/go-fundamentals/gameapp/config/httpservercfg/constant"
	"golang.project/go-fundamentals/gameapp/repository/migrator"
	"golang.project/go-fundamentals/gameapp/repository/mysql"
	"golang.project/go-fundamentals/gameapp/repository/redis/redismatching"
	"golang.project/go-fundamentals/gameapp/scheduler"
	"golang.project/go-fundamentals/gameapp/service/authenticationservice"
	"golang.project/go-fundamentals/gameapp/service/matchingservice"
	"log"
	"strings"
	"time"
)

type AppConfig struct {
	GracefullyShutdownTimeout time.Duration `mapstructure:"gracefully_shutdown_timeout"`
}

type HttpServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type Config struct {
	AppCfg                AppConfig                    `mapstructure:"app_cfg"`
	ServerCfg             HttpServerConfig             `mapstructure:"httpserver_cfg"`
	DataBaseCfg           mysql.Config                 `mapstructure:"database_cfg"`
	AuthCfg               authenticationservice.Config `mapstructure:"auth_cfg"`
	MatchingCfg           matchingservice.Config       `mapstructure:"matching_cfg"`
	RedisCfg              redis.Config                 `mapstructure:"redis_cfg"`
	SchedulerCfg          scheduler.Config             `mapstructure:"scheduler_cfg"`
	MatchingRepoCfg       redismatching.Config         `mapstructure:"matching_repo_cfg"`
	GrpcPresenceClientCfg presenceclient.Config        `mapstructure:"grpc_presence_client_cfg"`
}

func NewConfig(host string, port int) Config {

	cfg := loadConfig(host, port)

	return Config{
		AppCfg:                cfg.AppCfg,
		ServerCfg:             cfg.ServerCfg,
		DataBaseCfg:           cfg.DataBaseCfg,
		AuthCfg:               cfg.AuthCfg,
		MatchingCfg:           cfg.MatchingCfg,
		RedisCfg:              cfg.RedisCfg,
		SchedulerCfg:          cfg.SchedulerCfg,
		MatchingRepoCfg:       cfg.MatchingRepoCfg,
		GrpcPresenceClientCfg: cfg.GrpcPresenceClientCfg,
	}
}

// 1. read config file
// 2. env variable
func loadConfig(host string, port int) Config {

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// read config file
	viper.SetConfigName(constant.DefaultConfigFileName)
	viper.SetConfigType(constant.DefaultConfigFileType)
	viper.AddConfigPath(constant.DefaultConfigFilePath)

	var cfg Config
	if err := viper.ReadInConfig(); err != nil {

		log.Println("⚠️ config file not found, using environment variables")

		// get config from env variable
		if uErr := viper.Sub("httpserver_cfg").Unmarshal(&cfg.ServerCfg); uErr != nil {
			log.Fatalf("can't unmarshal httpserver config: %v", uErr)
		}
		if uErr := viper.Sub("database_cfg").Unmarshal(&cfg.DataBaseCfg); uErr != nil {
			log.Fatalf("can't unmarshal database config: %v", uErr)
		}
		if uErr := viper.Sub("auth_cfg").Unmarshal(&cfg.AuthCfg); uErr != nil {
			log.Fatalf("can't unmarshal auth config: %v", uErr)
		}
		if uErr := viper.Sub("matching_cfg").Unmarshal(&cfg.MatchingCfg); uErr != nil {
			log.Fatalf("can't unmarshal matching config: %v", uErr)
		}
		if uErr := viper.Sub("redis_cfg").Unmarshal(&cfg.RedisCfg); uErr != nil {
			log.Fatalf("can't unmarshal redis config: %v", uErr)
		}
		if uErr := viper.Sub("app_cfg").Unmarshal(&cfg.AppCfg); uErr != nil {
			log.Fatalf("can't unmarshal application config: %v", uErr)
		}
		if uErr := viper.Sub("scheduler_cfg").Unmarshal(&cfg.SchedulerCfg); uErr != nil {
			log.Fatalf("can't unmarshal scheduler config: %v", uErr)
		}
		if uErr := viper.Sub("matching_repo_cfg").Unmarshal(&cfg.MatchingRepoCfg); uErr != nil {
			log.Fatalf("can't unmarshal matching_repo_cfg config: %v", uErr)
		}
		if uErr := viper.Sub("grpc_presence_client_cfg").Unmarshal(&cfg.GrpcPresenceClientCfg); uErr != nil {
			log.Fatalf("can't unmarshal grpc_presence_client_cfg config: %v", uErr)
		}
	} else {

		if uErr := viper.Unmarshal(&cfg); uErr != nil {

			panic(fmt.Errorf("can't Unmarshal config file into struct Config, %w", uErr))
		}
	}

	if host != "" {
		cfg.ServerCfg.Host = host
	}
	if port != 0 {
		cfg.ServerCfg.Port = port
	}

	return cfg
}

func (c Config) Migrate(migrationCommand string) {

	if migrationCommand != "up" && migrationCommand != "down" && migrationCommand != "skip" && migrationCommand != "status" {
		panic(fmt.Sprintf("invalid migration-command: %s", migrationCommand))
	}

	dbConnection := mysql.NewDB(c.DataBaseCfg).MysqlConnection
	c.migrate(dbConnection, constant.MigrateDialect, migrationCommand)

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
