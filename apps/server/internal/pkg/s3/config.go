package s3

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Host        string `envconfig:"S3_HOST"`
	AccessToken string `envconfig:"S3_ACCESS_TOKEN"`
	SecretToken string `envconfig:"S3_SECRET_TOKEN"`
	Region      string `envconfig:"S3_REGION"`
	Bucket      string `envconfig:"S3_BUCKET"`
	Secure      bool   `envconfig:"S3_SECURE"`
}

func NewConfig() (Config, error) {
	cfg := Config{}

	wd, err := os.Getwd()
	if err != nil {
		return cfg, err
	}

	envPath := filepath.Join(wd, ".env")

	_ = godotenv.Load(envPath)

	if err := envconfig.Process("", &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
