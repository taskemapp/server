package app

import (
	"context"
	"fmt"
	"github.com/taskemapp/server/apps/server/internal/app/grpcfx"
	"github.com/taskemapp/server/apps/server/internal/app/profilefx"
	"github.com/taskemapp/server/apps/server/internal/logger"
	"net/url"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/taskemapp/server/apps/server/internal/app/authfx"
	"github.com/taskemapp/server/apps/server/internal/app/taskfx"
	"github.com/taskemapp/server/apps/server/internal/app/teamfx"
	"github.com/taskemapp/server/apps/server/internal/config"
	"github.com/taskemapp/server/apps/server/internal/pkg/migrations"
	"github.com/taskemapp/server/apps/server/internal/pkg/s3"
	"github.com/taskemapp/server/libs/queue"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

var App = fx.Options(
	fx.Provide(fx.Annotate(setupLogger, fx.As(new(logger.Logger)))),

	fx.Provide(setupConfig),
	fx.Provide(setupPgPool),
	fx.Provide(setupRabbitMq),
	fx.Provide(setupRedisClient),

	fx.Provide(queue.NewConfig),
	fx.Provide(fx.Annotate(queue.NewMQ, fx.As(new(queue.Queue)))),

	authfx.App,
	teamfx.App,
	profilefx.App,
	taskfx.App,
	grpcfx.App,

	fx.Invoke(
		migrations.Invoke,
	),
)

func setupConfig() (config.Config, error) {
	var cfg config.Config
	qCfg, err := queue.NewConfig()
	if err != nil {
		return cfg, err
	}
	s3Cfg, err := s3.NewConfig()
	if err != nil {
		return cfg, err
	}

	cfg, err = config.New(qCfg, s3Cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func setupLogger(c config.Config) (logger.Logger, error) {
	var zc zap.Config

	switch c.AppEnv {
	case envDev:
		zc = zap.NewDevelopmentConfig()
	case envProd:
		zc = zap.NewProductionConfig()
	default:
		zc = zap.NewDevelopmentConfig()
	}

	zc.OutputPaths = []string{"stdout"}
	zc.ErrorOutputPaths = []string{"stderr"}

	l, err := logger.New(&zc)
	if err != nil {
		return nil, err
	}
	return l, err
}

func setupPgPool(c config.Config) (*pgxpool.Pool, error) {
	return pgxpool.New(context.Background(), c.PostgresUrl)
}

func setupRabbitMq(c config.Config) (*amqp.Connection, error) {
	return amqp.Dial(c.RabbitMq.Url)
}

func setupRedisClient(c config.Config) (*redis.Client, error) {
	redisURL, err := url.Parse(c.RedisURL)
	if err != nil {
		return nil, err
	}

	addr := redisURL.Host

	password, _ := redisURL.User.Password()

	var db int
	if redisURL.Path != "" {
		_, err = fmt.Sscanf(redisURL.Path, "/%d", &db)
		if err != nil {
			return nil, err
		}
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}
