package template

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGet(t *testing.T) {
	tests := []struct {
		temp   Template
		failed bool
	}{
		{
			temp:   VerifyEmailTemplate,
			failed: false,
		},
	}

	for _, tt := range tests {
		_, err := Get(tt.temp, WithDir("../.."))
		assert.NoError(t, err)
	}
}
