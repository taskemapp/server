package notifier

type EmailNotifier interface {
	Send(msg EmailMsg) error
}

type EmailMsg struct {
	From    string
	To      []string
	Subject string
	Body    string
}

type PushNotifier interface {
	Send(msg PushMsg) error
}

type PushMsg struct {
	From     string
	To       []string
	Title    string
	Topic    string
	Subtitle string
	Body     string
	Image    string
	Priority int
}
