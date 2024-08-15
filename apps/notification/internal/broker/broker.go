package broker

type Broker interface {
	Receive(queueName string) error
}
