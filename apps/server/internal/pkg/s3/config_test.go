package s3

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	testCfg := Config{
		Host:        "http://minio:9001",
		AccessToken: "root",
		SecretToken: "password",
		Region:      "eu-west-1",
		Bucket:      "taskem",
		Secure:      false,
	}

	cfg, err := NewConfig()
	assert.NoError(t, err)

	assert.Equal(t, cfg, testCfg)
}
