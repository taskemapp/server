package broker

import (
	"time"
)

type Channel string

const (
	NotificationChannel Channel = "notification"
)

type Broker interface {
	Send(opts SendOpts, chName Channel) error
}

type SendOpts struct {
	ContentType string
	Headers     map[string]interface{}
	Timestamp   time.Time
	Body        []byte
}
