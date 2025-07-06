package publisher

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"golang.project/go-fundamentals/gameapp/adapter/redis"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/metrics"
	"time"
)

type Config struct {
	ContextTimeoutRedisPub time.Duration `mapstructure:"context_timeout_redis_pub"`
}
type Publisher struct {
	config  Config
	adapter *redis.Adapter
}

func NewPublisher(config Config, adapter *redis.Adapter) Publisher {
	return Publisher{
		config:  config,
		adapter: adapter,
	}
}

func (p Publisher) PublishEvent(event string, payload interface{}) {
	ctx, cancel := context.WithTimeout(context.Background(), p.config.ContextTimeoutRedisPub)
	defer cancel()
	err := p.adapter.GetClient().Publish(ctx, event, payload).Err()
	if err != nil {
		metrics.FailedPublishedEventCounter.Inc()
		logger.Warn(err, fmt.Sprintf("failed to publish event: %s", event))
	}

	metrics.PublishedEventCounter.With(prometheus.Labels{"event_name": event}).Inc()
}
