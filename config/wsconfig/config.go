package wsconfig

import (
	"github.com/spf13/viper"
	"golang.project/go-fundamentals/gameapp/adapter/presenceclient"
	"golang.project/go-fundamentals/gameapp/adapter/redis"
	"golang.project/go-fundamentals/gameapp/adapter/subscriber"
	"golang.project/go-fundamentals/gameapp/delivery/metricsserver"
	"golang.project/go-fundamentals/gameapp/gateway/websocket"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/pkg/jwt"
	"strings"
)

type Config struct {
	WSCfg                 websocket.Config      `mapstructure:"ws_cfg"`
	JwtCfg                jwt.Config            `mapstructure:"jwt_cfg"`
	LoggerCfg             logger.Config         `mapstructure:"logger_cfg"`
	GrpcPresenceClientCfg presenceclient.Config `mapstructure:"grpc_presence_client_cfg"`
	RedisCfg              redis.Config          `mapstructure:"redis_cfg"`
	SubscriberCfg         subscriber.Config     `mapstructure:"subscriber_cfg"`
	MetricsCfg            metricsserver.Config  `mapstructure:"websocket_metrics_cfg"`
}

func NewConfig() *Config {
	return &Config{
		JwtCfg:                jwt.Config{},
		WSCfg:                 websocket.Config{},
		LoggerCfg:             logger.Config{},
		GrpcPresenceClientCfg: presenceclient.Config{},
		RedisCfg:              redis.Config{},
		SubscriberCfg:         subscriber.Config{},
		MetricsCfg:            metricsserver.Config{},
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
		if uErr := viper.Sub("jwt_cfg").Unmarshal(&c.JwtCfg); uErr != nil {
			logger.Fatal(uErr, "can't unmarshal jwt config")
		}
		if uErr := viper.Sub("ws_cfg").Unmarshal(&c.WSCfg); uErr != nil {
			logger.Fatal(uErr, "can't unmarshal web socket config")
		}
		if uErr := viper.Sub("logger_cfg").Unmarshal(&c.LoggerCfg); uErr != nil {
			logger.Fatal(uErr, "can't unmarshal logger_cfg config")
		}
		if uErr := viper.Sub("grpc_presence_client_cfg").Unmarshal(&c.GrpcPresenceClientCfg); uErr != nil {
			logger.Fatal(uErr, "can't unmarshal grpc_presence_client_cfg config")
		}
		if uErr := viper.Sub("redis_cfg").Unmarshal(&c.RedisCfg); uErr != nil {
			logger.Fatal(uErr, "can't unmarshal redis config")
		}
		if uErr := viper.Sub("subscriber_cfg").Unmarshal(&c.SubscriberCfg); uErr != nil {
			logger.Fatal(uErr, "can't unmarshal subscriber config")
		}
		if uErr := viper.Sub("metrics_cfg").Unmarshal(&c.MetricsCfg); uErr != nil {
			logger.Fatal(uErr, "can't unmarshal metrics_cfg config")
		}
	} else {
		if uErr := viper.Unmarshal(c); uErr != nil {
			logger.Panic(uErr, "can't Unmarshal config file into struct Config")
		}
	}

	if host != "" {
		c.WSCfg.Host = host
	}
	if port != 0 {
		c.WSCfg.Port = port
	}

	return c
}
