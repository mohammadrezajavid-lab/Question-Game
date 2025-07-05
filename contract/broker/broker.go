package broker

import "context"

type Publisher interface {
	Published(event string, payload interface{})
}

type Subscriber interface {
	Subscribed(ctx context.Context, topic string) (<-chan string, error)
}
