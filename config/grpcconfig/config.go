package grpcconfig

import (
	"fmt"
	"github.com/spf13/viper"
	"golang.project/go-fundamentals/gameapp/adapter/redis"
	"golang.project/go-fundamentals/gameapp/service/presenceservice"
	"log"
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
}

func NewConfig() Config {
	return Config{
		GrpcCfg:     GrpcServerConfig{},
		RedisCfg:    redis.Config{},
		PresenceCfg: presenceservice.Config{},
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

		log.Println("⚠️ config file not found, using environment variables")

		// get config from env variable
		if uErr := viper.Sub("grpc_server_cfg").Unmarshal(&cfg.GrpcCfg); uErr != nil {
			log.Fatalf("can't unmarshal grpc server config: %v", uErr)
		}
		if uErr := viper.Sub("redis_cfg").Unmarshal(&cfg.RedisCfg); uErr != nil {
			log.Fatalf("can't unmarshal redis config: %v", uErr)
		}
		if uErr := viper.Sub("presence_cfg").Unmarshal(&cfg.PresenceCfg); uErr != nil {
			log.Fatalf("can't unmarshal presence config: %v", uErr)
		}

	} else {

		if uErr := viper.Unmarshal(&cfg); uErr != nil {
			panic(fmt.Errorf("can't Unmarshal config file into struct Config, %w", uErr))
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
