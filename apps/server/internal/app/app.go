package app

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/taskemapp/server/apps/server/internal/app/auth"
	grpcsrv "github.com/taskemapp/server/apps/server/internal/app/grpc"
	"github.com/taskemapp/server/apps/server/internal/app/task"
	"github.com/taskemapp/server/apps/server/internal/app/team"
	"github.com/taskemapp/server/apps/server/internal/broker"
	"github.com/taskemapp/server/apps/server/internal/config"
	"github.com/taskemapp/server/apps/server/internal/grpc/interceptors"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"net/url"
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
	fx.Provide(fx.Annotate(broker.New, fx.As(new(broker.Broker)))),

	auth.App,
	team.App,
	task.App,

	fx.Provide(interceptors.New),
	fx.Provide(grpcsrv.New),

	fx.Invoke(
		func(p *pgxpool.Pool, c config.Config, log *zap.Logger) error {
			if err := goose.SetDialect("pgx"); err != nil {
				log.Sugar().Error("Failed to set dialect: ", err)
				return err
			}
			db, err := sql.Open("pgx", c.PostgresUrl)
			if err != nil {
				log.Sugar().Error("Failed to open db conn: ", err)
				return err
			}
			defer db.Close()

			log.Sugar().Info("Run migrations")
			err = goose.Up(db, "migrations")
			if err != nil {
				log.Sugar().Error("Migration failed: ", err)
				return err
			}

			return nil
		},
		func(lc fx.Lifecycle, log *zap.Logger, c config.Config, srv *grpc.Server) {
			lc.Append(
				fx.Hook{
					OnStart: func(ctx context.Context) error {
						log.Sugar().Infof("Server starting on port %d", c.GrpcPort)

						l, err := net.Listen("tcp", fmt.Sprintf(":%d", c.GrpcPort))
						if err != nil {
							return err
						}

						reflection.Register(srv)

						go func() {
							err = srv.Serve(l)
							if err != nil {
								log.Error(err.Error())
								return
							}
						}()

						return nil
					},
					OnStop: func(ctx context.Context) error {
						log.Sugar().Info("Gracefully stopping grpc server")
						srv.GracefulStop()

						return nil
					},
				},
			)
		},
	),
)

func setupConfig() (config.Config, error) {
	cfg, err := config.New()
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
	return amqp.Dial(c.RabbitMqUrl)
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
