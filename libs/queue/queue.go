package broker

type ConsumeFn func(<-chan Message) error

type Queue interface {
	Publish(queue string, message Message) error
	Consume(name string, handler ConsumeFn) error
	Close() error
}

type Message struct {
	ContentType string
	Body        []byte
}
