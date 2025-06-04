package publisher

import (
	"context"
	"golang.project/go-fundamentals/gameapp/adapter/redis"
	"log"
	"time"
)

type Publish struct {
	redisAdapter *redis.Adapter
}

func NewPublish(redisAdapter *redis.Adapter) Publish {
	return Publish{redisAdapter: redisAdapter}
}

func (p Publish) PublishEvent(event string, payload interface{}) {
	const operation = "publisher.PublishedEvent"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := p.redisAdapter.GetClient().Publish(ctx, event, payload).Err()
	if err != nil {
		log.Printf("operation[%s], Error publishing Event: %v", operation, err.Error())
	}
}
