package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	AppEnv string `envconfig:"APP_ENV" default:"dev"`

	GrpcPort int `envconfig:"GRPC_PORT" default:"50051"`

	TokenTtl time.Duration `envconfig:"TOKEN_TTL" default:"1h"`
	// RefreshTokenTtl default value is 168 hours = 1 week
	RefreshTokenTtl time.Duration `envconfig:"REFRESH_TOKEN_TTL" default:"168h"`
	TokenSecret     string        `envconfig:"TOKEN_SECRET" required:"true"`

	PostgresUrl string `envconfig:"POSTGRES_URL" required:"true"`
	RedisURL    string `envconfig:"REDIS_URL" required:"true"`
	RabbitMqUrl string `envconfig:"RABBITMQ_URL" required:"true"`

	S3Host        string `envconfig:"S3_HOST"`
	S3AccessToken string `envconfig:"S3_ACCESS_TOKEN"`
	S3SecretToken string `envconfig:"S3_SECRET_TOKEN"`
	S3Region      string `envconfig:"S3_REGION"`
	S3Bucket      string `envconfig:"S3_BUCKET"`
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

	return cfg, nil
}
