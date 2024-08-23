package notifier

const (
	EmailNotify = "email"
	PushNotify  = "push"
)

type Notification struct {
	Type      string `json:"type"`
	Recipient string `json:"recipient"`
	Subject   string `json:"subject"`
	Message   string `json:"message"`
}
