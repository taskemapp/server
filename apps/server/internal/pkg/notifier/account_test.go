package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/taskemapp/server/apps/notification/pkg/notifier"
	"github.com/taskemapp/server/apps/server/internal/config"
	"github.com/taskemapp/server/libs/queue"
	"github.com/taskemapp/server/libs/template"
	"testing"
)

type testGenerator struct {
	cfg config.Config
}

func (t *testGenerator) VerifyLink() (link string, id uuid.UUID, err error) {
	return fmt.Sprintf("%s/verify?id=%s", t.cfg.HostDomain, uuid.Nil), uuid.Nil, nil
}

func (t *testGenerator) UnsubLink() (link string, id uuid.UUID, err error) {
	return fmt.Sprintf("%s/unsub?id=%s", t.cfg.HostDomain, uuid.Nil), uuid.Nil, nil
}

type args struct {
	name  string
	to    string
	title string
	from  string
}

func setupMock(t *testing.T, tt args, generator LinkGenerator) *queue.MockQueue {
	queueMock := queue.NewMockQueue(t)

	temp, err := template.Get(template.VerifyEmailTemplate)
	require.NoError(t, err)

	verifyLink, _, err := generator.VerifyLink()
	require.NoError(t, err)

	unsubLink, _, err := generator.UnsubLink()
	require.NoError(t, err)

	data := template.VerifyEmail{
		Name:             tt.name,
		ConfirmationLink: verifyLink,
		UnsubscribeLink:  unsubLink,
	}

	var buff bytes.Buffer
	err = temp.Execute(&buff, data)
	require.NoError(t, err)

	body, err := json.Marshal(notifier.EmailNotification{
		Notification: notifier.Notification{
			Title:   tt.title,
			Message: buff.String(),
		},
		To:   tt.to,
		From: tt.from,
	})

	message := queue.Message{
		ContentType: "application/json",
		Body:        body,
	}

	queueMock.EXPECT().Publish(notifier.ChannelEmail, message).Return(nil).Once()

	return queueMock
}

func TestEmailAccountNotifier_VerifyEmail(t *testing.T) {
	cfg := config.Config{}
	tt := args{
		name:  "ripls",
		to:    "ripls@taskem.test",
		title: "",
		from:  cfg.NoReplayEmail,
	}

	generator := testGenerator{cfg}

	mock := setupMock(t, tt, &generator)
	defer mock.AssertExpectations(t)

	opts := OptsEmailAccNotifier{
		Config: config.Config{},
		Queue:  mock,
		Lg:     &generator,
	}

	err := NewEmailAccountNotifier(opts).VerifyEmail(tt.name, tt.to)
	assert.NoError(t, err)
}
