package app

import (
	"context"
	"fmt"
	"github.com/taskemapp/server/apps/server/internal/pkg/notifier"
	"net/url"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/taskemapp/server/apps/server/internal/app/auth"
	"github.com/taskemapp/server/apps/server/internal/app/grpc"
	grpcsrv "github.com/taskemapp/server/apps/server/internal/app/grpc"
	"github.com/taskemapp/server/apps/server/internal/app/task"
	"github.com/taskemapp/server/apps/server/internal/app/team"
	"github.com/taskemapp/server/apps/server/internal/config"
	"github.com/taskemapp/server/apps/server/internal/grpc/interceptors"
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
	fx.Provide(setupConfig),
	fx.Provide(setupLogger),
	fx.Provide(setupPgPool),
	fx.Provide(setupRabbitMq),
	fx.Provide(setupRedisClient),
	fx.Provide(fx.Annotate(func(cfg config.Config) *notifier.BasicGenerator {
		return &notifier.BasicGenerator{HostDomain: cfg.HostDomain}
	}, fx.As(new(notifier.LinkGenerator)))),

	//RabbitMq
	fx.Provide(queue.NewConfig),
	fx.Provide(fx.Annotate(queue.NewMQ, fx.As(new(queue.Queue)))),

	//S3
	fx.Provide(s3.NewConfig),
	fx.Provide(s3.New),

	//General app
	auth.App,
	team.App,
	task.App,
	fx.Provide(interceptors.New),
	fx.Provide(grpcsrv.New),

	fx.Invoke(
		migrations.Invoke,
		s3.Invoke,
		grpc.Invoke,
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

func setupLogger(c config.Config) *zap.Logger {
	var log *zap.Logger

	switch c.AppEnv {
	case envDev:
		log, _ = zap.NewDevelopment()
	case envProd:
		log, _ = zap.NewProduction()
	default:
		log, _ = zap.NewDevelopment()
	}

	return log
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
		fmt.Sscanf(redisURL.Path, "/%d", &db)
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
