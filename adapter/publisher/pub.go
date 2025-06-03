package publisher

import (
	"context"
	"golang.project/go-fundamentals/gameapp/adapter/redis"
	"golang.project/go-fundamentals/gameapp/pkg/richerror"
	"log"
)

type Publish struct {
	redisAdapter *redis.Adapter
}

func NewPublish(redisAdapter *redis.Adapter) Publish {
	return Publish{redisAdapter: redisAdapter}
}

func (p Publish) PublishEvent(ctx context.Context, topic string, payload interface{}) error {
	const operation = "publisher.PublishedEvent"

	err := p.redisAdapter.GetClient().Publish(ctx, topic, payload).Err()
	if err != nil {
		log.Printf("Error publishing Event: %v", err.Error())

		return richerror.NewRichError(operation).WithError(err).WithKind(richerror.KindUnexpected)
	}

	return nil
}
