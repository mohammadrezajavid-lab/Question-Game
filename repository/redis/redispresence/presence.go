package redispresence

import (
	"context"
	"golang.project/go-fundamentals/gameapp/pkg/richerror"
	"time"
)

func (r *RedisDb) Upsert(ctx context.Context, key string, timestamp int64, expirationTime time.Duration) error {
	const operation = "redispresence.Upsert"

	if err := r.redisAdapter.GetClient().Set(ctx, key, timestamp, expirationTime).Err(); err != nil {
		return richerror.NewRichError(operation).WithError(err).WithKind(richerror.KindUnexpected)
	}

	return nil
}
