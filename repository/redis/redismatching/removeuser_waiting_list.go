package redismatching

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/metrics"
	"strconv"
)

func (r *RedisDb) RemoveUserFromWaitingList(userIds []uint, key string) {

	if len(userIds) < 1 {
		//logger.Info("No user IDs provided to remove. Skipping ZRem.")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.config.ContextTimeOutForZRem)
	defer cancel()

	members := make([]any, 0, len(userIds))
	for _, userId := range userIds {
		members = append(members, strconv.Itoa(int(userId)))
	}

	_, err := r.redisAdapter.GetClient().ZRem(ctx, key, members...).Result()
	if err != nil {
		metrics.RedisRequestsCounter.With(prometheus.Labels{"status": "fail"}).Inc()
		metrics.FailedZRemRedisCounter.Inc()

		logger.Info("failed_zrem_redis")
	}

	metrics.RedisRequestsCounter.With(prometheus.Labels{"status": "success"}).Inc()
}
