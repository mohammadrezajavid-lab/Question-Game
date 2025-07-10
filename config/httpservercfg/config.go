package httpservercfg

import (
	"database/sql"
	"errors"
	"github.com/spf13/viper"
	"golang.project/go-fundamentals/gameapp/adapter/presenceclient"
	"golang.project/go-fundamentals/gameapp/adapter/publisher"
	"golang.project/go-fundamentals/gameapp/adapter/quizclient"
	"golang.project/go-fundamentals/gameapp/adapter/redis"
	"golang.project/go-fundamentals/gameapp/adapter/subscriber"
	"golang.project/go-fundamentals/gameapp/config/httpservercfg/constant"
	"golang.project/go-fundamentals/gameapp/delivery/grpcserver/quizserver"
	"golang.project/go-fundamentals/gameapp/delivery/metricsserver"
	"golang.project/go-fundamentals/gameapp/delivery/pprofserver"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/pkg/jwt"
	"golang.project/go-fundamentals/gameapp/repository/migrator"
	"golang.project/go-fundamentals/gameapp/repository/mysql"
	"golang.project/go-fundamentals/gameapp/repository/redis/redismatching"
	"golang.project/go-fundamentals/gameapp/repository/redis/redisquiz"
	"golang.project/go-fundamentals/gameapp/scheduler"
	"golang.project/go-fundamentals/gameapp/service/gameservice"
	"golang.project/go-fundamentals/gameapp/service/matchingservice"
	"golang.project/go-fundamentals/gameapp/service/quizservice"
	"strings"
	"time"
)

type AppConfig struct {
	GracefullyShutdownTimeout time.Duration `mapstructure:"gracefully_shutdown_timeout"`
	DebugMod                  bool          `mapstructure:"debug_mod"`
}

type HttpServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type Config struct {
	AppCfg                AppConfig              `mapstructure:"app_cfg"`
	ServerCfg             HttpServerConfig       `mapstructure:"httpserver_cfg"`
	DataBaseCfg           mysql.Config           `mapstructure:"database_cfg"`
	JwtCfg                jwt.Config             `mapstructure:"jwt_cfg"`
	MatchingCfg           matchingservice.Config `mapstructure:"matching_cfg"`
	GameServiceCfg        gameservice.Config     `mapstructure:"game_svc_cfg"`
	QuizServiceCfg        quizservice.Config     `mapstructure:"quiz_svc_cfg"`
	RedisCfg              redis.Config           `mapstructure:"redis_cfg"`
	SchedulerCfg          scheduler.Config       `mapstructure:"scheduler_cfg"`
	MatchingRepoCfg       redismatching.Config   `mapstructure:"matching_repo_cfg"`
	QuizRedisRepoCfg      redisquiz.Config       `mapstructure:"quiz_redis_repo_cfg"`
	GrpcPresenceClientCfg presenceclient.Config  `mapstructure:"grpc_presence_client_cfg"`
	GrpcQuizClientCfg     quizclient.Config      `mapstructure:"grpc_quiz_client_cfg"`
	GrpcQuizCfg           quizserver.Config      `mapstructure:"grpc_quiz_server_cfg"`
	PublisherCfg          publisher.Config       `mapstructure:"publisher_cfg"`
	SubscriberCfg         subscriber.Config      `mapstructure:"subscriber_cfg"`
	LoggerCfg             logger.Config          `mapstructure:"logger_cfg"`
	MetricsCfg            metricsserver.Config   `mapstructure:"metrics_cfg"`
	PprofCfg              pprofserver.Config     `mapstructure:"pprof_cfg"`
}

func NewConfig(host string, port int) Config {

	cfg := loadConfig(host, port)

	return Config{
		AppCfg:                cfg.AppCfg,
		ServerCfg:             cfg.ServerCfg,
		DataBaseCfg:           cfg.DataBaseCfg,
		JwtCfg:                cfg.JwtCfg,
		MatchingCfg:           cfg.MatchingCfg,
		GameServiceCfg:        cfg.GameServiceCfg,
		QuizServiceCfg:        cfg.QuizServiceCfg,
		RedisCfg:              cfg.RedisCfg,
		SchedulerCfg:          cfg.SchedulerCfg,
		MatchingRepoCfg:       cfg.MatchingRepoCfg,
		QuizRedisRepoCfg:      cfg.QuizRedisRepoCfg,
		GrpcPresenceClientCfg: cfg.GrpcPresenceClientCfg,
		GrpcQuizClientCfg:     cfg.GrpcQuizClientCfg,
		GrpcQuizCfg:           cfg.GrpcQuizCfg,
		PublisherCfg:          cfg.PublisherCfg,
		SubscriberCfg:         cfg.SubscriberCfg,
		LoggerCfg:             cfg.LoggerCfg,
		MetricsCfg:            cfg.MetricsCfg,
		PprofCfg:              cfg.PprofCfg,
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
		logger.Info("config file not found, using environment variables")

		// get config from env variable
		if uErr := viper.Sub("httpserver_cfg").Unmarshal(&cfg.ServerCfg); uErr != nil {
			logger.Fatal(uErr, "can't unmarshal httpserver config")
		}
		if uErr := viper.Sub("database_cfg").Unmarshal(&cfg.DataBaseCfg); uErr != nil {
			logger.Fatal(uErr, "can't unmarshal database config")
		}
		if uErr := viper.Sub("jwt_cfg").Unmarshal(&cfg.JwtCfg); uErr != nil {
			logger.Fatal(uErr, "can't unmarshal JWT config")
		}
		if uErr := viper.Sub("matching_cfg").Unmarshal(&cfg.MatchingCfg); uErr != nil {
			logger.Fatal(uErr, "can't unmarshal matching config")
		}
		if uErr := viper.Sub("game_svc_cfg").Unmarshal(&cfg.GameServiceCfg); uErr != nil {
			logger.Fatal(uErr, "can't unmarshal game_svc_cfg config")
		}
		if uErr := viper.Sub("quiz_svc_cfg").Unmarshal(&cfg.QuizServiceCfg); uErr != nil {
			logger.Fatal(uErr, "can't unmarshal quiz_svc_cfg config")
		}
		if uErr := viper.Sub("redis_cfg").Unmarshal(&cfg.RedisCfg); uErr != nil {
			logger.Fatal(uErr, "can't unmarshal redis config")
		}
		if uErr := viper.Sub("app_cfg").Unmarshal(&cfg.AppCfg); uErr != nil {
			logger.Fatal(uErr, "can't unmarshal application config")
		}
		if uErr := viper.Sub("scheduler_cfg").Unmarshal(&cfg.SchedulerCfg); uErr != nil {
			logger.Fatal(uErr, "can't unmarshal scheduler config")
		}
		if uErr := viper.Sub("matching_repo_cfg").Unmarshal(&cfg.MatchingRepoCfg); uErr != nil {
			logger.Fatal(uErr, "can't unmarshal matching_repo_cfg config")
		}
		if uErr := viper.Sub("quiz_redis_repo_cfg").Unmarshal(&cfg.QuizRedisRepoCfg); uErr != nil {
			logger.Fatal(uErr, "can't unmarshal quiz_redis_repo_cfg config")
		}
		if uErr := viper.Sub("grpc_presence_client_cfg").Unmarshal(&cfg.GrpcPresenceClientCfg); uErr != nil {
			logger.Fatal(uErr, "can't unmarshal grpc_presence_client_cfg config")
		}
		if uErr := viper.Sub("grpc_quiz_client_cfg").Unmarshal(&cfg.GrpcQuizClientCfg); uErr != nil {
			logger.Fatal(uErr, "can't unmarshal grpc_quiz_client_cfg config")
		}
		if uErr := viper.Sub("grpc_quiz_server_cfg").Unmarshal(&cfg.GrpcQuizCfg); uErr != nil {
			logger.Fatal(uErr, "can't unmarshal grpc_quiz_server_cfg config")
		}
		if uErr := viper.Sub("publisher_cfg").Unmarshal(&cfg.PublisherCfg); uErr != nil {
			logger.Fatal(uErr, "can't unmarshal publisher_cfg config")
		}
		if uErr := viper.Sub("subscriber_cfg").Unmarshal(&cfg.SubscriberCfg); uErr != nil {
			logger.Fatal(uErr, "can't unmarshal subscriber_cfg config")
		}
		if uErr := viper.Sub("logger_cfg").Unmarshal(&cfg.LoggerCfg); uErr != nil {
			logger.Fatal(uErr, "can't unmarshal logger_cfg config")
		}
		if uErr := viper.Sub("metrics_cfg").Unmarshal(&cfg.MetricsCfg); uErr != nil {
			logger.Fatal(uErr, "can't unmarshal metrics_cfg config")
		}
		if uErr := viper.Sub("pprof_cfg").Unmarshal(&cfg.PprofCfg); uErr != nil {
			logger.Fatal(uErr, "can't unmarshal pprof_cfg config")
		}
	} else {

		if uErr := viper.Unmarshal(&cfg); uErr != nil {
			logger.Panic(uErr, "can't Unmarshal config file into struct Config")
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
		logger.Warn(errors.New("invalid migration-command, use default [skip]"), "invalid migration-command")
		migrationCommand = "skip"
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
