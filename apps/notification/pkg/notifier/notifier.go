package notifier

const (
	ChannelEmail = "email.notify"
	ChannelPush  = "push.notify"
)

type Notification struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type EmailNotification struct {
	Notification
	To   string `json:"to"`
	From string `json:"from"`
}

type PushNotification struct {
	Notification
	Tokens   []string               `json:"tokens"`
	Topic    string                 `json:"topic,omitempty"`
	Image    string                 `json:"image,omitempty"`
	Subtitle string                 `json:"subtitle,omitempty"`
	Priority int                    `json:"priority,omitempty"`
	Data     map[string]interface{} `json:"data,omitempty"`
}
