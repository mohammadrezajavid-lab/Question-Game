package redismatching

import (
	"context"
	"golang.project/go-fundamentals/gameapp/entity"
	"log"
	"strconv"
)

func (r *RedisDb) RemoveUserFromWaitingList(userIds []uint, category entity.Category) {

	const operation = "redismatching.RemoveUserFromWaitedList"

	if len(userIds) < 1 {
		log.Printf("%s: No user IDs provided to remove. Skipping ZRem.", operation)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.config.ContextTimeOutForZRem)
	defer cancel()

	var key = r.GetKey(category)

	members := make([]any, 0, len(userIds))
	for _, userId := range userIds {
		members = append(members, strconv.Itoa(int(userId)))
	}

	numberOfRemovedMember, err := r.redisAdapter.GetClient().ZRem(ctx, key, members...).Result()
	if err != nil {
		// TODO - update metrics
		// TODO - log error
		log.Printf("%s: Error ZRem from waiting list: %v\n", operation, err)
	}

	// TODO - update metrics
	// TODO - log error
	log.Printf("%d items removed from %s", numberOfRemovedMember, key)
}
