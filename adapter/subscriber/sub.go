package subscriber

import (
	"context"
	"golang.project/go-fundamentals/gameapp/adapter/redis"
)

type Config struct {
	QueueBufferSize uint `mapstructure:"queue_buffer_size"`
}

type Subscriber struct {
	adapter *redis.Adapter
	config  Config
}

func NewSubscriber(adapter *redis.Adapter, config Config) Subscriber {
	return Subscriber{
		adapter: adapter,
		config:  config,
	}
}

func (s Subscriber) SubscribeTopic(ctx context.Context, topic string) (<-chan string, error) {
	sub := s.adapter.GetClient().Subscribe(ctx, topic)
	ch := make(chan string, s.config.QueueBufferSize)

	go func() {
		defer sub.Close()
		defer close(ch)

		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-sub.Channel():
				if !ok {
					return
				}
				ch <- msg.Payload
			}
		}
	}()

	return ch, nil
}
