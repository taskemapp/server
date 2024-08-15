package broker

import (
	"github.com/google/uuid"
	"time"
)

type Broker interface {
	Send(opts SendOpts, chName string) error
}

type SendOpts struct {
	ContentType string
	Headers     map[string]interface{}
	Timestamp   time.Time
	UserId      uuid.UUID
	Body        []byte
}
