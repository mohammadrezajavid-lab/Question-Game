package broker

import "context"

type Publisher interface {
	PublishEvent(event string, payload interface{})
}

type Subscriber interface {
	SubscribeTopic(ctx context.Context, topic string) (<-chan string, error)
}
