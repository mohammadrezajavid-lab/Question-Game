package quizgrpcconfig

import (
	"github.com/spf13/viper"
	"golang.project/go-fundamentals/gameapp/adapter/redis"
	"golang.project/go-fundamentals/gameapp/delivery/grpcserver/quizserver"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/repository/mysql"
	"golang.project/go-fundamentals/gameapp/repository/redis/redisset"
	"golang.project/go-fundamentals/gameapp/service/quizservice"
	"strings"
)

type Config struct {
	QuizServiceCfg   quizservice.Config `mapstructure:"quiz_svc_cfg"`
	GrpcQuizCfg      quizserver.Config  `mapstructure:"grpc_quiz_server_cfg"`
	QuizRedisRepoCfg redisset.Config    `mapstructure:"quiz_redis_repo_cfg"`
	DataBaseCfg      mysql.Config       `mapstructure:"database_cfg"`
	RedisCfg         redis.Config       `mapstructure:"redis_cfg"`
	LoggerCfg        logger.Config      `mapstructure:"logger_cfg"`
}

func NewConfig() *Config {
	return &Config{
		QuizServiceCfg:   quizservice.Config{},
		GrpcQuizCfg:      quizserver.Config{},
		QuizRedisRepoCfg: redisset.Config{},
		DataBaseCfg:      mysql.Config{},
		RedisCfg:         redis.Config{},
		LoggerCfg:        logger.Config{},
	}
}

func (c *Config) LoadConfig(host string, port int) *Config {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// read config file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		logger.Warn(err, "config file not found, using environment variables")

		// get config from env variable
		if uErr := viper.Sub("quiz_svc_cfg").Unmarshal(&c.QuizServiceCfg); uErr != nil {
			logger.Fatal(uErr, "can't unmarshal quiz_svc_cfg config")
		}
		if uErr := viper.Sub("grpc_quiz_server_cfg").Unmarshal(&c.GrpcQuizCfg); uErr != nil {
			logger.Fatal(uErr, "can't unmarshal grpc_quiz_server_cfg config")
		}
		if uErr := viper.Sub("quiz_redis_repo_cfg").Unmarshal(&c.QuizRedisRepoCfg); uErr != nil {
			logger.Fatal(uErr, "can't unmarshal quiz_redis_repo_cfg config")
		}
		if uErr := viper.Sub("database_cfg").Unmarshal(&c.DataBaseCfg); uErr != nil {
			logger.Fatal(uErr, "can't unmarshal database_cfg config")
		}
		if uErr := viper.Sub("redis_cfg").Unmarshal(&c.RedisCfg); uErr != nil {
			logger.Fatal(uErr, "can't unmarshal redis_cfg config")
		}
		if uErr := viper.Sub("logger_cfg").Unmarshal(&c.LoggerCfg); uErr != nil {
			logger.Fatal(uErr, "can't unmarshal logger_cfg config")
		}
	} else {
		if uErr := viper.Unmarshal(c); uErr != nil {
			logger.Panic(uErr, "can't Unmarshal config file into struct Config")
		}
	}

	if host != "" {
		c.GrpcQuizCfg.Host = host
	}
	if port != 0 {
		c.GrpcQuizCfg.Port = port
	}

	return c
}
