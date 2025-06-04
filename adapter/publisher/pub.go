package publisher

import (
	"context"
	"golang.project/go-fundamentals/gameapp/adapter/redis"
	"log"
	"time"
)

type Config struct {
	ContextTimeoutRedisPub time.Duration `mapstructure:"context_timeout_redis_pub"`
}
type Publish struct {
	config       Config
	redisAdapter *redis.Adapter
}

func NewPublish(config Config, redisAdapter *redis.Adapter) Publish {
	return Publish{
		config:       config,
		redisAdapter: redisAdapter,
	}
}

func (p Publish) PublishEvent(event string, payload interface{}) {
	const operation = "publisher.PublishedEvent"

	ctx, cancel := context.WithTimeout(context.Background(), p.config.ContextTimeoutRedisPub)
	defer cancel()
	err := p.redisAdapter.GetClient().Publish(ctx, event, payload).Err()
	if err != nil {
		// TODO - update metrics
		// TODO - log error
		log.Printf("operation[%s], Error publishing Event: %v", operation, err.Error())
	}

	// TODO - update metrics
	// TODO - log error
}
