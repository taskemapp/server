package auth

import (
	authserver "github.com/taskemapp/server/apps/server/internal/grpc/auth"
	"github.com/taskemapp/server/apps/server/internal/pkg/notifier"
	"github.com/taskemapp/server/apps/server/internal/repository/token"
	"github.com/taskemapp/server/apps/server/internal/repository/user"
	authservice "github.com/taskemapp/server/apps/server/internal/service/auth"
	"go.uber.org/fx"
)

var App = fx.Options(
	fx.Provide(
		fx.Annotate(user.NewPgx, fx.As(new(user.Repository))),
	),

	fx.Provide(
		fx.Annotate(token.NewClient, fx.As(new(token.Repository))),
	),

	fx.Provide(
		fx.Annotate(notifier.NewEmailAccountNotifier, fx.As(new(notifier.AccountNotifier))),
	),

	fx.Provide(
		fx.Annotate(authservice.New, fx.As(new(authservice.Service))),
	),
	fx.Provide(authserver.New),
)
