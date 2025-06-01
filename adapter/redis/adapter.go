package redis

import (
	"fmt"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Network  string `mapstructure:"network_cfg"` // "tcp" or "unix"
	Host     string `mapstructure:"host_cfg"`
	Port     int    `mapstructure:"port_cfg"`
	Password string `mapstructure:"password_cfg"`
	DB       int    `mapstructure:"db_cfg"`
}

type Adapter struct {
	client *redis.Client
}

func New(cfg *Config) *Adapter {
	return &Adapter{client: redis.NewClient(&redis.Options{
		Network:  cfg.Network,
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})}
}

func (a *Adapter) GetClient() *redis.Client {
	return a.client
}
