package notifier

import (
	"github.com/go-faster/errors"
	"github.com/taskemapp/server/apps/server/internal/config"
	"github.com/taskemapp/server/libs/queue"
	"github.com/taskemapp/server/libs/template"
	"go.uber.org/fx"
)

type AccountNotifier interface {
	VerifyEmail(username string, email string) error
}

type EmailAccountNotifier struct {
	config config.Config
	q      queue.Queue
	lg     LinkGenerator
}

type OptsEmailAccNotifier struct {
	fx.In
	Config config.Config
	Queue  queue.Queue
	Lg     LinkGenerator
}

func NewEmailAccountNotifier(opts OptsEmailAccNotifier) *EmailAccountNotifier {
	return &EmailAccountNotifier{
		config: opts.Config,
		q:      opts.Queue,
		lg:     opts.Lg,
	}
}

func (n *EmailAccountNotifier) VerifyEmail(username string, email string) error {

	confirmLink, _, err := n.lg.VerifyLink()
	if err != nil {
		return errors.Wrap(err, "verify email")
	}

	unsubLink, _, err := n.lg.UnsubLink()
	if err != nil {
		return errors.Wrap(err, "verify email")
	}

	temp, err := template.Get(template.VerifyEmailTemplate)
	if err != nil {
		return errors.Wrap(err, "verify email")
	}

	err = n.sendEmail(sendEmailOpts{
		temp: temp,
		data: template.VerifyEmail{
			Name:             username,
			ConfirmationLink: confirmLink,
			UnsubscribeLink:  unsubLink,
		},
		title: "",
		to:    email,
		from:  n.config.NoReplayEmail,
	})
	if err != nil {
		return errors.Wrap(err, "verify email")
	}

	return nil
}
