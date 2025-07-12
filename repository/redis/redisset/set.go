package redisset

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/metrics"
	"time"
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

func (r *RedisDb) Set(ctx context.Context, key string, value string, ttlExpiration time.Duration) error {
	err := r.redisAdapter.GetClient().Set(ctx, key, value, ttlExpiration).Err()
	if err != nil {
		metrics.RedisRequestsCounter.With(prometheus.Labels{"status": "fail"}).Inc()
		logger.Warn(err, fmt.Sprintf("failed to Set key: %s", key))
		return err
	}
	metrics.RedisRequestsCounter.With(prometheus.Labels{"status": "success"}).Inc()
	return nil
}

func (r *RedisDb) Get(ctx context.Context, key string) (string, error) {
	value, err := r.redisAdapter.GetClient().Get(ctx, key).Result()
	if err != nil {
		metrics.RedisRequestsCounter.With(prometheus.Labels{"status": "fail"}).Inc()
		logger.Warn(err, fmt.Sprintf("failed to Get key: %s", key))
		return "", err
	}
	metrics.RedisRequestsCounter.With(prometheus.Labels{"status": "success"}).Inc()
	return value, nil
}

func (r *RedisDb) SetAdd(ctx context.Context, key string, value string) {
	err := r.redisAdapter.GetClient().SAdd(ctx, key, value).Err()
	if err != nil {
		metrics.RedisRequestsCounter.With(prometheus.Labels{"status": "fail"}).Inc()
		logger.Warn(err, fmt.Sprintf("failed to SAdd key: %s", key))
		return
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
