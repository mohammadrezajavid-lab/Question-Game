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

//GetKey
/*
* Naming convention,
* waiting_list_prefix:{difficulty}:{category}
 */
func (r *RedisDb) GetKey(category entity.Category, difficulty entity.QuestionDifficulty) string {
	return fmt.Sprintf("%s:%s:%d", r.config.WaitingListPrefix, category, difficulty)
}
