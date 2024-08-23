package notifier

import (
	"github.com/wneessen/go-mail"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Opts struct {
	fx.In
	*zap.Logger
	*mail.Client
}

type Email struct {
	log *zap.Logger
	c   *mail.Client
}

func NewEmail(opts Opts) *Email {
	return &Email{
		log: opts.Logger,
		c:   opts.Client,
	}
}

func (e *Email) Send(opts EmailOpts) error {
	e.log.Info("Creating email")

	m := mail.NewMsg()
	if err := m.From(opts.From); err != nil {
		e.log.Sugar().Errorf("Failed to create email: %s", err)
		return err
	}
	if err := m.To(opts.To...); err != nil {
		e.log.Sugar().Errorf("Failed to create email: %s", err)
		return err
	}

	m.Subject(opts.Subject)
	m.SetBodyString(mail.TypeTextHTML, "Do you like this mail? I certainly do!")

	e.log.Info("Sending email")

	err := e.c.DialAndSend(m)
	if err != nil {
		e.log.Sugar().Errorf("Failed to send email: %s", err)
	}

	e.log.Info("Email sent")

	return nil
}
