package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	AppEnv string `envconfig:"APP_ENV" default:"dev"`

	PostgresUrl string `envconfig:"POSTGRES_URL"`
	RedisURL    string `envconfig:"REDIS_URL"`

	S3Host        string `envconfig:"S3_HOST"`
	S3AccessToken string `envconfig:"S3_ACCESS_TOKEN"`
	S3SecretToken string `envconfig:"S3_SECRET_TOKEN"`
	S3Region      string `envconfig:"S3_REGION"`
	S3Bucket      string `envconfig:"S3_BUCKET"`
}

func New() (Config, error) {
	cfg := Config{}

	err := godotenv.Load()

	if err != nil {
		return cfg, err
	}

	if err := envconfig.Process("", &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
