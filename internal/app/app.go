package app

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"os"
	"path/filepath"
	"taskem-server/internal/config"
	grpcsrv "taskem-server/internal/grpc"
	authserver "taskem-server/internal/grpc/auth"
	"taskem-server/internal/repositories/user"
	authservice "taskem-server/internal/service/auth"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

var App = fx.Options(
	fx.Provide(setupConfig),
	fx.Provide(setupLogger),
	fx.Provide(setupPgPool),

	fx.Provide(
		fx.Annotate(user.NewPgx, fx.As(new(user.Repository))),
	),

	fx.Provide(
		fx.Annotate(authservice.New, fx.As(new(authservice.Auth))),
	),
	fx.Provide(authserver.New),
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
	}

	return log
}

func setupPgPool(c config.Config) (*pgxpool.Pool, error) {
	return pgxpool.New(context.Background(), c.PostgresUrl)
}
