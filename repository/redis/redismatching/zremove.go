package redismatching

import (
	"context"
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/pkg/richerror"
	"strconv"
)

func (r *RedisDb) RemoveUserFromWaitedList(ctx context.Context, userId uint, category entity.Category) error {

	const operation = "redismatching.RemoveUserFromWaitedList"

	var key = r.GetKey(category)
	_, err := r.redisAdapter.GetClient().ZRem(ctx, key, strconv.Itoa(int(userId))).Result()
	if err != nil {
		return richerror.NewRichError(operation).
			WithError(err).
			WithKind(richerror.KindUnexpected).
			WithMeta(map[string]interface{}{"userId": userId, "category": category})
	}

	return nil
}
