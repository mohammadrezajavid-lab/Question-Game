package redispresence

import (
	"context"
	"github.com/redis/go-redis/v9"
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/pkg/richerror"
	"strconv"
	"strings"
	"time"
)

func (r *RedisDb) Upsert(ctx context.Context, key string, timestamp int64, expirationTime time.Duration) error {
	const operation = "redispresence.Upsert"

	if err := r.redisAdapter.GetClient().Set(ctx, key, timestamp, expirationTime).Err(); err != nil {
		return richerror.NewRichError(operation).WithError(err).WithKind(richerror.KindUnexpected)
	}

	return nil
}

//func (r *RedisDb) GetPresence(ctx context.Context, key string) (entity.Presence, error) {
//	const operation = "redispresence.GetPresence"
//
//	timestamp, err := r.redisAdapter.GetClient().Get(ctx, key).Result()
//	if err != nil {
//		if errors.Is(err, redis.Nil) {
//			return entity.Presence{}, richerror.NewRichError(operation).WithError(err).WithMessage("presence not found").WithKind(richerror.KindNotFound)
//		}
//
//		return entity.Presence{}, richerror.NewRichError(operation).WithError(err).WithKind(richerror.KindUnexpected)
//	}
//
//	ts, aErr := strconv.Atoi(timestamp)
//	if aErr != nil {
//		return entity.Presence{}, richerror.NewRichError(operation).WithError(aErr).WithKind(richerror.KindUnexpected)
//	}
//
//	return entity.NewPresence(0, int64(ts)), nil
//}

func (r *RedisDb) GetPresences(ctx context.Context, keys []string) ([]entity.Presence, error) {
	const operation = "redispresence.GetPresences"
	if len(keys) == 0 {
		return []entity.Presence{}, nil
	}

	values, err := r.redisAdapter.GetClient().MGet(ctx, keys...).Result()
	if err != nil {
		return nil, richerror.NewRichError(operation).WithError(err).WithKind(richerror.KindUnexpected)
	}

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
