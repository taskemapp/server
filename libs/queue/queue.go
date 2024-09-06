package queue

type ConsumeFn func(msg Message)

type Queue interface {
	Publish(queue string, message Message) error
	Consume(name string, handler ConsumeFn) error
	Close() error
}

type Message struct {
	ContentType string
	Body        []byte
}
