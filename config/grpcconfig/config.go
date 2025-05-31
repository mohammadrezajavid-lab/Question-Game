package grpcconfig

import (
	"fmt"
	"github.com/spf13/viper"
	"golang.project/go-fundamentals/gameapp/adapter/redis"
	"golang.project/go-fundamentals/gameapp/config/httpservercfg/constant"
	"golang.project/go-fundamentals/gameapp/service/presenceservice"
	"log"
	"strings"
)

type GrpcServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type Config struct {
	GrpcCfg     GrpcServerConfig       `mapstructure:"grpc_server_cfg"`
	RedisCfg    redis.Config           `mapstructure:"redis_cfg"`
	PresenceCfg presenceservice.Config `mapstructure:"presence_cfg"`
}

func NewConfig(host string, port int) Config {
	cfg := loadConfig(host, port)

	return Config{
		GrpcCfg:     cfg.GrpcCfg,
		RedisCfg:    cfg.RedisCfg,
		PresenceCfg: cfg.PresenceCfg,
	}
}

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
