package redismatching

import "golang.project/go-fundamentals/gameapp/adapter/redis"

type Config struct {
	WaitingListPrefix string `mapstructure:"waiting_list_prefix"`
}

type RedisDb struct {
	config       Config
	redisAdapter *redis.Adapter
}

func NewRedisDb(redisAdapter *redis.Adapter) *RedisDb {
	return &RedisDb{redisAdapter: redisAdapter}
}
