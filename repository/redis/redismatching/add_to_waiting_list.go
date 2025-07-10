package redismatching

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"golang.project/go-fundamentals/gameapp/metrics"
	"golang.project/go-fundamentals/gameapp/pkg/errormessage"
	"golang.project/go-fundamentals/gameapp/pkg/richerror"
	"golang.project/go-fundamentals/gameapp/pkg/timestamp"
	"strconv"
)

func (r *RedisDb) AddToWaitingList(ctx context.Context, userId uint, key string) error {

	const operation = "redismatching.AddToWaitingList"

	rdb := r.redisAdapter.GetClient()

	timeStamp := timestamp.Now()

	// score: timestamp and member: userId
	member := redis.Z{
		Score:  float64(timeStamp),
		Member: strconv.Itoa(int(userId)),
	}

	if _, aErr := rdb.ZAdd(ctx, key, member).Result(); aErr != nil {
		metrics.RedisRequestsCounter.With(prometheus.Labels{"status": "fail"}).Inc()

		return richerror.NewRichError(operation).
			WithError(aErr).
			WithMessage(errormessage.ErrorMsgUnexpected).
			WithKind(richerror.KindUnexpected)
	}

	metrics.RedisRequestsCounter.With(prometheus.Labels{"status": "success"}).Inc()

	return nil
}
