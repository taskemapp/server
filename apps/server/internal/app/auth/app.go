package auth

import (
	"go.uber.org/fx"
	authserver "server/internal/grpc/auth"
	"server/internal/repositories/user"
	authservice "server/internal/service/auth"
)

var App = fx.Options(
	fx.Provide(
		fx.Annotate(user.NewPgx, fx.As(new(user.Repository))),
	),

	fx.Provide(
		fx.Annotate(authservice.New, fx.As(new(authservice.Service))),
	),
	fx.Provide(authserver.New),
)
