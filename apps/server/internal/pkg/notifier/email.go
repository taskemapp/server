package notifier

import (
	"bytes"
	"encoding/json"
	"html/template"

	"github.com/go-faster/errors"
	"github.com/taskemapp/server/apps/notification/pkg/notifier"
	"github.com/taskemapp/server/libs/queue"
)

type sendEmailOpts struct {
	temp  *template.Template
	data  any
	title string
	to    string
	from  string
}

func (n *EmailAccountNotifier) sendEmail(opts sendEmailOpts) error {
	var buff bytes.Buffer
	err := opts.temp.Execute(&buff, opts.data)
	if err != nil {
		return errors.Wrap(err, "send email")
	}

	body, err := json.Marshal(notifier.EmailNotification{
		Notification: notifier.Notification{
			Title:   opts.title,
			Message: buff.String(),
		},
		To:   opts.to,
		From: opts.from,
	})
	if err != nil {
		return errors.Wrap(err, "send email")
	}

	msg := queue.Message{
		ContentType: "application/json",
		Body:        body,
	}

	err = n.q.Publish(notifier.ChannelEmail, msg)
	if err != nil {
		return errors.Wrap(err, "send email")
	}

	return nil
}
