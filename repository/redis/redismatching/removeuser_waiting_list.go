package redismatching

import (
	"context"
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/metrics"
	"strconv"
)

func (r *RedisDb) RemoveUserFromWaitingList(userIds []uint, category entity.Category) {

	const operation = "redismatching.RemoveUserFromWaitedList"

	if len(userIds) < 1 {
		//logger.Info("No user IDs provided to remove. Skipping ZRem.")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.config.ContextTimeOutForZRem)
	defer cancel()

	var key = r.GetKey(category)

	members := make([]any, 0, len(userIds))
	for _, userId := range userIds {
		members = append(members, strconv.Itoa(int(userId)))
	}

	_, err := r.redisAdapter.GetClient().ZRem(ctx, key, members...).Result()
	if err != nil {
		metrics.FailedZRemRedisCounter.Inc()
		logger.Info("failed_zrem_redis")
	}
}
