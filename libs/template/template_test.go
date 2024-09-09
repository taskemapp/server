package template

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGet(t *testing.T) {
	tests := []struct {
		temp   Type
		failed bool
	}{
		{
			temp:   VerifyEmailTemplate,
			failed: false,
		},
	}

	for _, tt := range tests {
		temp, err := Get(tt.temp, WithDir("../.."))
		assert.NoError(t, err)

		switch tt.temp {
		case VerifyEmailTemplate:
			var buff bytes.Buffer
			err = temp.Execute(&buff, VerifyEmail{
				Name:             "name",
				ConfirmationLink: "conf-link",
				UnsubscribeLink:  "unsubscribe-link",
			})
			assert.NoError(t, err)
		}

	}
}
