package wsconfig

import (
	"github.com/spf13/viper"
	"golang.project/go-fundamentals/gameapp/adapter/presenceclient"
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
}

func NewConfig() *Config {
	return &Config{
		JwtCfg:                jwt.Config{},
		WSCfg:                 websocket.Config{},
		LoggerCfg:             logger.Config{},
		GrpcPresenceClientCfg: presenceclient.Config{},
	}
}

func (c *Config) LoadConfig(host string, port int) *Config {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// read config file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	//var cfg Config
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
