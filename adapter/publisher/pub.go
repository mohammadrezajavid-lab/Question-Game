package publisher

import (
	"context"
	"golang.project/go-fundamentals/gameapp/adapter/redis"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/metrics"
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

	ctx, cancel := context.WithTimeout(context.Background(), p.config.ContextTimeoutRedisPub)
	defer cancel()
	err := p.redisAdapter.GetClient().Publish(ctx, event, payload).Err()
	if err != nil {
		metrics.FailedPublishedEventCounter.Inc()
		logger.Warn(err, "failed_published_event")
	}

	metrics.PublishedEventCounter.Inc()
}
