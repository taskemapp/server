package notifier

import (
	"github.com/wneessen/go-mail"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type PushOpts struct {
	fx.In
	*zap.Logger
	*mail.Client
}

type Fcm struct {
	log *zap.Logger
	c   *mail.Client
}

func NewFcm(opts PushOpts) *Fcm {
	return &Fcm{
		log: opts.Logger,
		c:   opts.Client,
	}
}

func (e *Fcm) Send(opts PushMsg) error {

	e.log.Info("Fcm sent")

	return nil
}
