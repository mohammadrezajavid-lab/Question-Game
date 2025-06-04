package broker

type Published interface {
	PublishEvent(event string, payload interface{})
}
