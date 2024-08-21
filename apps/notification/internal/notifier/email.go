package notifier

import (
	"github.com/jordan-wright/email"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"sync"
	"time"
)

type Opts struct {
	fx.In
	*zap.Logger
	*email.Pool
}

type Email struct {
	log *zap.Logger
	c   *email.Pool
	wg  sync.WaitGroup
}

func NewEmail(opts Opts) *Email {
	return &Email{
		log: opts.Logger,
		c:   opts.Pool,
		wg:  sync.WaitGroup{},
	}
}

func (e *Email) Send() error {
	e.log.Info("Creating email")
	msg := email.NewEmail()

	msg.From = "Jordan Wright <kanada.smirnov@ya.ru>"
	msg.To = []string{"kanada.smirnov@gmail.com"}

	msg.Subject = "Awesome Subject"
	msg.Text = []byte("Text Body is, of course, supported!")
	msg.HTML = []byte("<h1>Fancy HTML is supported, too!</h1>")

	e.log.Info("Sending email")
	e.wg.Add(1)
	go func() {
		defer e.c.Close()

		err := e.c.Send(msg, 5*time.Second)
		if err != nil {
			e.log.Sugar().Error("Failed sending email: ", err)
		}

		e.wg.Done()
		e.log.Info("Email sent")
	}()

	return nil
}
