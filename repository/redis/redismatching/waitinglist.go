package redismatching

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/pkg/errormessage"
	"golang.project/go-fundamentals/gameapp/pkg/richerror"
	"strconv"
	"time"
)

func (r *RedisDb) AddToWaitingList(userId uint, category entity.Category) error {

	const operation = "redismatching.AddToWaitingList"

	rdb := r.redisAdapter.GetClient()
	ctx := context.Background()

	var key string = fmt.Sprintf("%s:%v", "category", category)
	timeStamp := time.Now().UnixNano()
	number := redis.Z{
		Score:  float64(timeStamp),
		Member: strconv.Itoa(int(userId)),
	}

	if _, aErr := rdb.ZAdd(ctx, key, number).Result(); aErr != nil {

		return richerror.NewRichError(operation).
			WithError(aErr).
			WithMessage(errormessage.ErrorMsgUnexpected).
			WithKind(richerror.KindUnexpected)
	}

	return nil
}
