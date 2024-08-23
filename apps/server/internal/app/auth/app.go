package auth

import (
	authserver "github.com/taskemapp/server/apps/server/internal/grpc/auth"
	"github.com/taskemapp/server/apps/server/internal/repositories/user"
	authservice "github.com/taskemapp/server/apps/server/internal/service/auth"
	"go.uber.org/fx"
	authserver "taskem-server/internal/grpc/auth"
	"taskem-server/internal/repositories/token"
	"taskem-server/internal/repositories/user"
	authservice "taskem-server/internal/service/auth"
)

var App = fx.Options(
	fx.Provide(
		fx.Annotate(user.NewPgx, fx.As(new(user.Repository))),
	),

	fx.Provide(
		fx.Annotate(token.NewClient, fx.As(new(token.Repository))),
	),

	fx.Provide(
		fx.Annotate(authservice.New, fx.As(new(authservice.Service))),
	),
	fx.Provide(authserver.New),
)
