package redismatching

import (
	"context"
	"golang.project/go-fundamentals/gameapp/entity"
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

	waitedUsers, zErr := r.redisAdapter.GetClient().ZRangeWithScores(ctx, key, 0, numRecords).Result()
	if zErr != nil {
		return make([]matchingparam.WaitedUser, 0),
			richerror.NewRichError(operation).WithError(zErr).WithKind(richerror.KindUnexpected)
	}

	waitedUsersList := make([]matchingparam.WaitedUser, 0, numRecords)
	for _, z := range waitedUsers {
		userId, _ := strconv.Atoi(z.Member.(string))
		waitedUsersList = append(waitedUsersList, matchingparam.NewWaitedUser(int64(z.Score), uint(userId), category))
	}

	return waitedUsersList, nil
}
