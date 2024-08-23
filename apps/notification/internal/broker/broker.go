package broker

type Channel string

const (
	NotificationChannel Channel = "notification"
)

type Broker interface {
	Receive(queueName string) error
}
