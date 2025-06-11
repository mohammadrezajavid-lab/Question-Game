package redismatching

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/metrics"
	"golang.project/go-fundamentals/gameapp/param/matchingparam"
	"golang.project/go-fundamentals/gameapp/pkg/richerror"
	"strconv"
)

func (r *RedisDb) GetWaitedUserByCategory(ctx context.Context, category entity.Category) ([]matchingparam.WaitedUser, error) {
	const operation = "redismatching.PopMinWaitedUserByCategory"

	var key = r.GetKey(category)
	var numRecords = r.redisAdapter.GetClient().ZCard(ctx, key).Val()
	if numRecords > 5000 {
		numRecords = numRecords / 4
	}
	metrics.RedisRequestsCounter.With(prometheus.Labels{"status": "success"}).Inc()

	waitedUsers, zErr := r.redisAdapter.GetClient().ZRangeWithScores(ctx, key, 0, numRecords).Result()
	if zErr != nil {
		metrics.RedisRequestsCounter.With(prometheus.Labels{"status": "fail"}).Inc()

		return make([]matchingparam.WaitedUser, 0),
			richerror.NewRichError(operation).WithError(zErr).WithKind(richerror.KindUnexpected)
	}
	metrics.RedisRequestsCounter.With(prometheus.Labels{"status": "success"}).Inc()

	waitedUsersList := make([]matchingparam.WaitedUser, 0, numRecords)
	for _, z := range waitedUsers {
		userId, _ := strconv.Atoi(z.Member.(string))
		waitedUsersList = append(waitedUsersList, matchingparam.NewWaitedUser(int64(z.Score), uint(userId), category))
	}

	return waitedUsersList, nil
}
