package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/taskemapp/server/apps/server/internal/pkg/s3"
	"github.com/taskemapp/server/libs/queue"
)

func TestNew(t *testing.T) {
	qTestCfg := queue.Config{
		Url: "amqp://taskem:password@0.0.0.0:56721",
	}
	qCfg, err := queue.NewConfig()
	assert.NoError(t, err)
	assert.Equal(t, qCfg, qTestCfg)

	s3TestCfg := s3.Config{
		Host:        "http://minio:9001",
		AccessToken: "root",
		SecretToken: "password",
		Region:      "eu-west-1",
		Bucket:      "taskem",
	}
	s3Cfg, err := s3.NewConfig()
	assert.NoError(t, err)
	assert.Equal(t, s3Cfg, s3TestCfg)

	_, err = New(qCfg, s3Cfg)
	assert.NoError(t, err)
}
