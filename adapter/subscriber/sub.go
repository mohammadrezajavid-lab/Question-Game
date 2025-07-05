package subscriber

import (
	"context"
	"golang.project/go-fundamentals/gameapp/adapter/redis"
)

type Config struct {
}
type Subscriber struct {
	config  Config
	adapter *redis.Adapter
}

func NewSubscriber(config Config, adapter *redis.Adapter) Subscriber {
	return Subscriber{
		config:  config,
		adapter: adapter,
	}
}

func (s Subscriber) Subscribed(ctx context.Context, topic string) (<-chan string, error) {

	sub := s.adapter.GetClient().Subscribe(ctx, topic)

	ch := make(chan string)

	go func() {
		defer sub.Close()
		for msg := range sub.Channel() {
			ch <- msg.Payload
		}
		close(ch)
	}()

	return ch, nil
}
