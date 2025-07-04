package redismatching

import (
	"fmt"
	"golang.project/go-fundamentals/gameapp/adapter/redis"
	"golang.project/go-fundamentals/gameapp/entity"
	"time"
)

type Config struct {
	WaitingListPrefix     string        `mapstructure:"waiting_list_prefix"`
	ContextTimeOutForZRem time.Duration `mapstructure:"context_timeout_ZRem"`
}

type RedisDb struct {
	config       Config
	redisAdapter *redis.Adapter
}

func NewRedisDb(redisAdapter *redis.Adapter, config Config) *RedisDb {
	return &RedisDb{
		config:       config,
		redisAdapter: redisAdapter,
	}
}

func (r *RedisDb) GetKey(category entity.Category) string {
	return fmt.Sprintf("%s:%s", r.config.WaitingListPrefix, category)
}
