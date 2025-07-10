package redisquiz

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"golang.project/go-fundamentals/gameapp/metrics"
)

func (r *RedisDb) SetLength(ctx context.Context, key string) (int, error) {
	length, err := r.redisAdapter.GetClient().SCard(ctx, key).Result()
	if err != nil {
		metrics.RedisRequestsCounter.With(prometheus.Labels{"status": "fail"}).Inc()
		return 0, err
	}

	metrics.RedisRequestsCounter.With(prometheus.Labels{"status": "success"}).Inc()
	return int(length), err
}

func (r *RedisDb) SetAdd(ctx context.Context, key string, value string) {
	err := r.redisAdapter.GetClient().SAdd(ctx, key, value).Err()
	if err != nil {
		metrics.RedisRequestsCounter.With(prometheus.Labels{"status": "fail"}).Inc()
	}
	metrics.RedisRequestsCounter.With(prometheus.Labels{"status": "success"}).Inc()
}

func (r *RedisDb) SetPop(ctx context.Context, key string) (string, error) {
	value, err := r.redisAdapter.GetClient().SPop(ctx, key).Result()
	if err != nil {
		metrics.RedisRequestsCounter.With(prometheus.Labels{"status": "fail"}).Inc()
		return "", err
	}

	metrics.RedisRequestsCounter.With(prometheus.Labels{"status": "success"}).Inc()
	return value, err
}
