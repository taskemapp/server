package notifier

type Notifier interface {
	Send() error
}
