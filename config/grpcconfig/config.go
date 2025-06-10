package grpcconfig

import (
	"github.com/spf13/viper"
	"golang.project/go-fundamentals/gameapp/adapter/redis"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/service/presenceservice"
	"strings"
)

type GrpcServerConfig struct {
	Host    string `mapstructure:"host"`
	Port    int    `mapstructure:"port"`
	Network string `mapstructure:"network"`
}

type Config struct {
	GrpcCfg     GrpcServerConfig       `mapstructure:"grpc_server_cfg"`
	RedisCfg    redis.Config           `mapstructure:"redis_cfg"`
	PresenceCfg presenceservice.Config `mapstructure:"presence_cfg"`
	LoggerCfg   logger.Config          `mapstructure:"logger_cfg"`
}

func NewConfig() Config {
	return Config{
		GrpcCfg:     GrpcServerConfig{},
		RedisCfg:    redis.Config{},
		PresenceCfg: presenceservice.Config{},
		LoggerCfg:   logger.Config{},
	}
}

func (c Config) LoadConfig(host string, port int) Config {

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// read config file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	var cfg Config
	if err := viper.ReadInConfig(); err != nil {

		logger.Info("config file not found, using environment variables")

		// get config from env variable
		if uErr := viper.Sub("grpc_server_cfg").Unmarshal(&cfg.GrpcCfg); uErr != nil {
			logger.Fatal(uErr, "can't unmarshal grpc server config")
		}
		if uErr := viper.Sub("redis_cfg").Unmarshal(&cfg.RedisCfg); uErr != nil {
			logger.Fatal(uErr, "can't unmarshal redis config")
		}
		if uErr := viper.Sub("presence_cfg").Unmarshal(&cfg.PresenceCfg); uErr != nil {
			logger.Fatal(uErr, "can't unmarshal presence config")
		}
		if uErr := viper.Sub("logger_cfg").Unmarshal(&cfg.LoggerCfg); uErr != nil {
			logger.Fatal(uErr, "can't unmarshal logger_cfg config")
		}

	} else {

		if uErr := viper.Unmarshal(&cfg); uErr != nil {
			logger.Panic(uErr, "can't Unmarshal config file into struct Config")
		}
	}

	if host != "" {
		cfg.GrpcCfg.Host = host
	}
	if port != 0 {
		cfg.GrpcCfg.Port = port
	}

	return cfg
}
