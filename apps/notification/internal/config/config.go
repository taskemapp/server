package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/taskemapp/server/libs/queue"
	"os"
	"path/filepath"
)

type Config struct {
	AppEnv string `envconfig:"APP_ENV" default:"dev"`

	RabbitMq queue.Config

	SmtpUrl      string `envconfig:"SMTP_URL" required:"true"`
	SmtpHost     string `envconfig:"SMTP_HOST" required:"true"`
	SmtpPort     string `envconfig:"SMTP_PORT" required:"true"`
	SmtpUsername string `envconfig:"SMTP_USERNAME" required:"true"`
	SmtpPassword string `envconfig:"SMTP_PASSWORD" required:"true"`
}

func New() (Config, error) {
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

	mqConfig, err := queue.NewConfig()
	if err != nil {
		panic(err)
	}

	cfg.RabbitMq = mqConfig

	return cfg, nil
}
