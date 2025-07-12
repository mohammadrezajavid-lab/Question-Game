package gameservicecfg

import (
	"github.com/spf13/viper"
	"golang.project/go-fundamentals/gameapp/adapter/publisher"
	"golang.project/go-fundamentals/gameapp/adapter/quizclient"
	"golang.project/go-fundamentals/gameapp/adapter/redis"
	"golang.project/go-fundamentals/gameapp/adapter/subscriber"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/repository/mysql"
	"golang.project/go-fundamentals/gameapp/service/gameservice"
	"strings"
)

type Config struct {
	GameServiceCfg    gameservice.Config `mapstructure:"game_svc_cfg"`
	RedisCfg          redis.Config       `mapstructure:"redis_cfg"`
	DataBaseCfg       mysql.Config       `mapstructure:"database_cfg"`
	PublisherCfg      publisher.Config   `mapstructure:"publisher_cfg"`
	SubscriberCfg     subscriber.Config  `mapstructure:"subscriber_cfg"`
	LoggerCfg         logger.Config      `mapstructure:"logger_cfg"`
	GrpcQuizClientCfg quizclient.Config  `mapstructure:"grpc_quiz_client_cfg"`
}

func NewConfig() *Config {
	return &Config{
		GameServiceCfg:    gameservice.Config{},
		RedisCfg:          redis.Config{},
		DataBaseCfg:       mysql.Config{},
		PublisherCfg:      publisher.Config{},
		SubscriberCfg:     subscriber.Config{},
		LoggerCfg:         logger.Config{},
		GrpcQuizClientCfg: quizclient.Config{},
	}
}

func (c *Config) LoadConfig() *Config {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// read config file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		logger.Warn(err, "config file not found, using environment variables")

		// get config from env variable
		if uErr := viper.Sub("game_svc_cfg").Unmarshal(&c.GameServiceCfg); uErr != nil {
			logger.Fatal(uErr, "can't unmarshal game service config")
		}
		if uErr := viper.Sub("redis_cfg").Unmarshal(&c.RedisCfg); uErr != nil {
			logger.Fatal(uErr, "can't unmarshal redis config")
		}
		if uErr := viper.Sub("database_cfg").Unmarshal(&c.DataBaseCfg); uErr != nil {
			logger.Fatal(uErr, "can't unmarshal database config")
		}
		if uErr := viper.Sub("publisher_cfg").Unmarshal(&c.PublisherCfg); uErr != nil {
			logger.Fatal(uErr, "can't unmarshal publisher config")
		}
		if uErr := viper.Sub("subscriber_cfg").Unmarshal(&c.SubscriberCfg); uErr != nil {
			logger.Fatal(uErr, "can't unmarshal subscriber config")
		}
		if uErr := viper.Sub("logger_cfg").Unmarshal(&c.LoggerCfg); uErr != nil {
			logger.Fatal(uErr, "can't unmarshal logger_cfg config")
		}
		if uErr := viper.Sub("grpc_quiz_client_cfg").Unmarshal(&c.GrpcQuizClientCfg); uErr != nil {
			logger.Fatal(uErr, "can't unmarshal grpc_quiz_client_cfg config")
		}
	} else {
		if uErr := viper.Unmarshal(c); uErr != nil {
			logger.Panic(uErr, "can't Unmarshal config file into struct Config")
		}
	}

	return c
}
