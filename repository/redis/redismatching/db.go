package redismatching

import "golang.project/go-fundamentals/gameapp/adapter/redis"

type RedisDb struct {
	redisAdapter *redis.Adapter
}

func NewRedisDb(redisAdapter *redis.Adapter) *RedisDb {
	return &RedisDb{redisAdapter: redisAdapter}
}
