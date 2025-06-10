package redismatching

import (
	"context"
	"github.com/redis/go-redis/v9"
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/pkg/errormessage"
	"golang.project/go-fundamentals/gameapp/pkg/richerror"
	"golang.project/go-fundamentals/gameapp/pkg/timestamp"
	"strconv"
)

func (r *RedisDb) AddToWaitingList(ctx context.Context, userId uint, category entity.Category) error {

	const operation = "redismatching.AddToWaitingList"

	rdb := r.redisAdapter.GetClient()

	var key = r.GetKey(category)
	timeStamp := timestamp.Now()

	// score: timestamp and member: userId
	member := redis.Z{
		Score:  float64(timeStamp),
		Member: strconv.Itoa(int(userId)),
	}

	if _, aErr := rdb.ZAdd(ctx, key, member).Result(); aErr != nil {

		return richerror.NewRichError(operation).
			WithError(aErr).
			WithMessage(errormessage.ErrorMsgUnexpected).
			WithKind(richerror.KindUnexpected)
	}

	return nil
}
