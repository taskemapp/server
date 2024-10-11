package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/taskemapp/server/apps/server/internal/pkg/s3"
	"github.com/taskemapp/server/libs/queue"
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

	RabbitMq queue.Config

	NoReplayEmail string `envconfig:"NO_REPLAY_EMAIL"`

	HostDomain string `envconfig:"HOST_DOMAIN"`

	S3 s3.Config
}

func New(
	qCfg queue.Config,
	s3Cfg s3.Config,
) (Config, error) {
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

	cfg.RabbitMq = qCfg
	cfg.S3 = s3Cfg

	return cfg, nil
}
