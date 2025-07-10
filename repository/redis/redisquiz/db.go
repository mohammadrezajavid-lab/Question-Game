package redisquiz

import (
	"fmt"
	"golang.project/go-fundamentals/gameapp/adapter/redis"
	"golang.project/go-fundamentals/gameapp/entity"
)

type Config struct {
	Prefix string `mapstructure:"prefix"`
}

type RedisDb struct {
	config       Config
	redisAdapter *redis.Adapter
}

func NewRedisDb(redisAdapter *redis.Adapter, config Config) *RedisDb {
	return &RedisDb{config: config, redisAdapter: redisAdapter}
}

//GetKey
/*
* Naming convention for key in key-value data structure,
* quiz-pool:{difficulty}:{category}
 */
func (r *RedisDb) GetKey(category entity.Category, difficulty entity.QuestionDifficulty) string {
	return fmt.Sprintf("%s:%d:%s", r.config.Prefix, difficulty, category)
}
