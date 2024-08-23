package notifier

type EmailNotifier interface {
	Send(opts EmailOpts) error
}

type EmailOpts struct {
	From    string
	To      []string
	Subject string
	Body    []byte
}

type PushNotifier interface {
	Send(opts PushOpts) error
}

type PushOpts struct {
	From    string
	To      []string
	Subject string
	Body    string
}
