package authfx

import (
	"github.com/taskemapp/server/apps/server/internal/config"
	authsrv "github.com/taskemapp/server/apps/server/internal/grpc/auth"
	"github.com/taskemapp/server/apps/server/internal/pkg/notifier"
	"github.com/taskemapp/server/apps/server/internal/repository/token"
	"github.com/taskemapp/server/apps/server/internal/repository/user"
	authservice "github.com/taskemapp/server/apps/server/internal/service/auth"
	v1 "github.com/taskemapp/server/apps/server/tools/gen/grpc/v1"
	"go.uber.org/fx"
)

var App = fx.Options(
	fx.Module(
		"auth",
		fx.Provide(
			fx.Private,
			fx.Annotate(
				func(cfg config.Config) *notifier.BasicGenerator {
					return &notifier.BasicGenerator{HostDomain: cfg.HostDomain}
				},
				fx.As(new(notifier.LinkGenerator)),
			),
			fx.Annotate(user.NewPgx, fx.As(new(user.Repository))),
			fx.Annotate(token.NewClient, fx.As(new(token.Repository))),
			fx.Annotate(notifier.NewEmailAccountNotifier, fx.As(new(notifier.AccountNotifier))),
			fx.Annotate(authservice.New, fx.As(new(authservice.Service))),
		),

		fx.Provide(
			fx.Annotate(authsrv.New, fx.As(new(v1.AuthServer))),
		),
	),
)
