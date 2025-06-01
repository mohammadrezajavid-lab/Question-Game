package redispresence

import (
	"context"
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/pkg/richerror"
	"strconv"
	"time"
)

func (r *RedisDb) Upsert(ctx context.Context, key string, timestamp int64, expirationTime time.Duration) error {
	const operation = "redispresence.Upsert"

	if err := r.redisAdapter.GetClient().Set(ctx, key, timestamp, expirationTime).Err(); err != nil {
		return richerror.NewRichError(operation).WithError(err).WithKind(richerror.KindUnexpected)
	}

	return nil
}

func (r *RedisDb) GetPresence(ctx context.Context, key string) (entity.Presence, error) {
	const operation = "redispresence.GetPresence"

	timestamp, err := r.redisAdapter.GetClient().Get(ctx, key).Result()
	if err != nil {
		return entity.Presence{}, richerror.NewRichError(operation).WithError(err).WithKind(richerror.KindUnexpected)
	}

	ts, aErr := strconv.Atoi(timestamp)
	if aErr != nil {
		return entity.Presence{}, richerror.NewRichError(operation).WithError(aErr).WithKind(richerror.KindUnexpected)
	}

	return entity.NewPresence(0, int64(ts)), nil
}
