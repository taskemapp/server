package notifier

import (
	"github.com/wneessen/go-mail"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type EmailOpts struct {
	fx.In
	*zap.Logger
	*mail.Client
}

type Email struct {
	log *zap.Logger
	c   *mail.Client
}

func NewEmail(opts EmailOpts) *Email {
	return &Email{
		log: opts.Logger,
		c:   opts.Client,
	}
}

func (e *Email) Send(opts EmailMsg) error {
	e.log.Debug("Creating email")

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

	m.SetBodyString(mail.TypeTextHTML, opts.Body)

	e.log.Debug("Sending email")

	err := e.c.DialAndSend(m)
	if err != nil {
		e.log.Sugar().Errorf("Failed to send email: %s", err)
		return err
	}

	e.log.Sugar().Infof("Email sent to: %s", opts.To)

	return nil
}
