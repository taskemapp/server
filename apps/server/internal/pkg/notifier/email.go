package notifier

import (
	"bytes"
	"encoding/json"
	"github.com/go-faster/errors"
	"github.com/taskemapp/server/apps/notification/pkg/notifier"
	"github.com/taskemapp/server/libs/queue"
	"html/template"
)

type sendEmailOpts struct {
	temp  *template.Template
	data  any
	q     queue.Queue
	title string
	to    string
	from  string
}

func sendEmail(opts sendEmailOpts) error {
	var buff bytes.Buffer
	err := opts.temp.Execute(&buff, opts.data)
	if err != nil {
		return errors.Wrap(err, "send email")
	}

	body, err := json.Marshal(notifier.EmailNotification{
		Notification: notifier.Notification{
			Title:   "Verify email",
			Message: buff.String(),
		},
		To:   opts.to,
		From: opts.from,
	})
	if err != nil {
		return errors.Wrap(err, "send email")
	}

	err = opts.q.Publish(notifier.ChannelEmail, queue.Message{
		ContentType: "application/json",
		Body:        body,
	})
	if err != nil {
		return errors.Wrap(err, "send email")
	}

	return nil
}
