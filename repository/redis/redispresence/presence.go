package redispresence

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/metrics"
	"golang.project/go-fundamentals/gameapp/pkg/richerror"
	"strconv"
	"strings"
	"time"
)

func (r *RedisDb) Upsert(ctx context.Context, key string, timestamp int64, expirationTime time.Duration) error {
	const operation = "redispresence.Upsert"

	if err := r.redisAdapter.GetClient().Set(ctx, key, timestamp, expirationTime).Err(); err != nil {
		metrics.RedisRequestsCounter.With(prometheus.Labels{"status": "fail"}).Inc()

		return richerror.NewRichError(operation).WithError(err).WithKind(richerror.KindUnexpected)
	}

	metrics.RedisRequestsCounter.With(prometheus.Labels{"status": "success"}).Inc()

	return nil
}

func (r *RedisDb) GetPresences(ctx context.Context, keys []string) ([]entity.Presence, error) {
	const operation = "redispresence.GetPresences"
	if len(keys) == 0 {
		return []entity.Presence{}, nil
	}

	values, err := r.redisAdapter.GetClient().MGet(ctx, keys...).Result()
	if err != nil {
		metrics.RedisRequestsCounter.With(prometheus.Labels{"status": "fail"}).Inc()

		return nil, richerror.NewRichError(operation).WithError(err).WithKind(richerror.KindUnexpected)
	}
	metrics.RedisRequestsCounter.With(prometheus.Labels{"status": "success"}).Inc()

	presences := make([]entity.Presence, 0, len(keys))
	for index, key := range keys {
		if values[index] == redis.Nil {
			continue
		}

		timestamp, _ := values[index].(string)
		timestampInt, _ := strconv.Atoi(timestamp)

		userId := strings.TrimPrefix(key, "presence:")
		userIdInt, _ := strconv.Atoi(userId)

		presences = append(presences, entity.NewPresence(uint(userIdInt), int64(timestampInt)))
	}

	return presences, nil
}
