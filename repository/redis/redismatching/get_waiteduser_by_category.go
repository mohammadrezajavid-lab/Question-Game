package redismatching

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/metrics"
	"golang.project/go-fundamentals/gameapp/param/gameparam"
	"golang.project/go-fundamentals/gameapp/pkg/richerror"
	"strconv"
)

func (r *RedisDb) GetWaitedUsersByCategory(ctx context.Context, category entity.Category, difficulty entity.QuestionDifficulty) ([]gameparam.WaitedUser, error) {
	const operation = "redismatching.PopMinWaitedUserByCategory"

	var key = r.GetKey(category, difficulty)
	var numRecords = r.redisAdapter.GetClient().ZCard(ctx, key).Val()
	if numRecords > 5000 {
		numRecords = numRecords / 4
	}
	metrics.RedisRequestsCounter.With(prometheus.Labels{"status": "success"}).Inc()

	waitedUsers, zErr := r.redisAdapter.GetClient().ZRangeWithScores(ctx, key, 0, numRecords).Result()
	if zErr != nil {
		metrics.RedisRequestsCounter.With(prometheus.Labels{"status": "fail"}).Inc()

		return make([]gameparam.WaitedUser, 0),
			richerror.NewRichError(operation).WithError(zErr).WithKind(richerror.KindUnexpected)
	}
	metrics.RedisRequestsCounter.With(prometheus.Labels{"status": "success"}).Inc()

	waitedUsersList := make([]gameparam.WaitedUser, 0, numRecords)
	for _, z := range waitedUsers {
		userId, _ := strconv.Atoi(z.Member.(string))
		waitedUsersList = append(waitedUsersList, gameparam.NewWaitedUser(int64(z.Score), uint(userId), category, difficulty))
	}

	return waitedUsersList, nil
}
