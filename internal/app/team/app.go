package team

import (
	"go.uber.org/fx"
	teamserver "taskem-server/internal/grpc/team"
	"taskem-server/internal/repositories/team"
)

var App = fx.Options(
	fx.Provide(
		fx.Annotate(team.NewPgx, fx.As(new(team.Repository))),
	),

	fx.Provide(teamserver.New),
)
