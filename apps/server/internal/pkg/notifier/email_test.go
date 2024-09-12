package notifier

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/taskemapp/server/apps/notification/pkg/notifier"
	"github.com/taskemapp/server/libs/queue"
	"github.com/taskemapp/server/libs/template"
	"testing"
)

func Test_sendEmail(t *testing.T) {
	type args struct {
		name  string
		to    string
		title string
		from  string
	}
	tt := args{
		name:  "",
		to:    "",
		title: "",
		from:  "",
	}

	//Setup mock
	queueMock := queue.NewMockQueue(t)
	defer queueMock.AssertExpectations(t)

	temp, err := template.Get(template.VerifyEmailTemplate, template.WithDir("../../../../.."))
	assert.NoError(t, err)

	data := template.VerifyEmail{
		Name:             tt.name,
		ConfirmationLink: "",
		UnsubscribeLink:  "",
	}

	var buff bytes.Buffer
	err = temp.Execute(&buff, data)
	assert.NoError(t, err)

	body, err := json.Marshal(notifier.EmailNotification{
		Notification: notifier.Notification{
			Title:   "Verify email",
			Message: buff.String(),
		},
		To:   tt.to,
		From: tt.from,
	})

	message := queue.Message{
		ContentType: "application/json",
		Body:        body,
	}

	queueMock.EXPECT().Publish(notifier.ChannelEmail, message).Return(nil)

	//Test starting
	err = queueMock.Publish(notifier.ChannelEmail, message)
	assert.NoError(t, err)

	err = sendEmail(sendEmailOpts{
		temp:  temp,
		data:  data,
		q:     queueMock,
		title: tt.title,
		to:    tt.to,
		from:  tt.from,
	})
	assert.NoError(t, err)
}
