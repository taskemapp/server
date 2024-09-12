package notifier

import (
	"github.com/go-faster/errors"
	"github.com/google/uuid"
	"github.com/taskemapp/server/apps/server/internal/config"
	"github.com/taskemapp/server/apps/server/internal/repositories/user"
	"github.com/taskemapp/server/libs/queue"
	"github.com/taskemapp/server/libs/template"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"strings"
)

type AccountNotifier interface {
	VerifyEmail(username string, email string) error
}

type EmailAccountNotifier struct {
	config   config.Config
	q        queue.Queue
	userRepo user.Repository
}

type OptsEmailAccNotifier struct {
	fx.In
	Config   config.Config
	Queue    queue.Queue
	UserRepo user.Repository
}

func NewEmailAccountNotifier(opts OptsEmailAccNotifier) *EmailAccountNotifier {
	return &EmailAccountNotifier{
		config:   opts.Config,
		q:        opts.Queue,
		userRepo: opts.UserRepo,
	}
}

func (n *EmailAccountNotifier) VerifyEmail(username string, email string) error {
	confirmID, err := uuid.NewUUID()
	if err != nil {
		return errors.Wrap(err, "verify to")
	}

	confirmLink, err := buildConfirmLink(n.config.HostDomain, confirmID.String())
	if err != nil {
		return errors.Wrap(err, "verify to")
	}

	temp, err := template.Get(template.VerifyEmailTemplate)
	if err != nil {
		return errors.Wrap(err, "verify to")
	}

	err = sendEmail(sendEmailOpts{
		temp: temp,
		data: template.VerifyEmail{
			Name:             username,
			ConfirmationLink: confirmLink,
			UnsubscribeLink:  "unsubscribe-link",
		},
		q:     n.q,
		title: "Verify to",
		to:    email,
		from:  n.config.NoReplayEmail,
	})
	if err != nil {
		return errors.Wrap(err, "verify to")
	}

	return nil
}

func buildConfirmLink(host string, confirmID string) (string, error) {
	var sb strings.Builder
	var err error

	_, err = sb.WriteString(host)
	err = multierr.Append(err, err)

	_, err = sb.WriteString("/verify?id=")
	err = multierr.Append(err, err)

	_, err = sb.WriteString(confirmID)
	err = multierr.Append(err, err)

	return sb.String(), err
}
