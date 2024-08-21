package app

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"taskem-server/internal/app/auth"
	grpcsrv "taskem-server/internal/app/grpc"
	"taskem-server/internal/app/task"
	"taskem-server/internal/app/team"
	"taskem-server/internal/config"
	"taskem-server/internal/grpc/interceptors"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

var App = fx.Options(
	fx.Provide(setupConfig),
	fx.Provide(setupLogger),
	fx.Provide(setupPgPool),
	fx.Provide(setupRedisClient),

	auth.App,
	team.App,
	task.App,

	fx.Provide(interceptors.New),
	fx.Provide(grpcsrv.New),

	fx.Invoke(
		func(p *pgxpool.Pool, c config.Config) error {
			if err := goose.SetDialect("pgx"); err != nil {
				return err
			}
			db, err := sql.Open("pgx", c.PostgresUrl)
			if err != nil {
				return err
			}
			defer db.Close()

			wd, err := os.Getwd()
			if err != nil {
				return err
			}

			return goose.Up(db, filepath.Join(wd, "migrations"))
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
