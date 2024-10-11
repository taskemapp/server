package s3

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	f, err := os.Create(".env")
	if err != nil {
		t.Fatal(err)
	}

	host := "http://minio:9001"
	access := "root"
	secret := "password"
	region := "eu-west-1"
	bucket := "taskem"
	secure := false

	testEnv := []byte(
		`# This is test data
		S3_ACCESS_TOKEN=` + access + `
		S3_SECRET_TOKEN=` + secret + `
		S3_REGION=` + region + `
		S3_HOST=` + host + `
		S3_BUCKET=` + bucket + `
		S3_SECURE=` + strconv.FormatBool(secure))
	f.Write(testEnv)
	f.Close()

	defer os.Remove(f.Name())

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

	assert.Equal(t, testCfg, cfg)
}
