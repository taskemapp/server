package team

import (
	"go.uber.org/fx"
	teamserver "taskem-server/internal/grpc/team"
	"taskem-server/internal/repositories/team"
	"taskem-server/internal/repositories/team_member"
	teamservice "taskem-server/internal/service/team"
)

var App = fx.Options(
	fx.Provide(
		fx.Annotate(team.NewPgx, fx.As(new(team.Repository))),
		fx.Annotate(team_member.NewPgx, fx.As(new(team_member.Repository))),
		fx.Annotate(teamservice.New, fx.As(new(teamservice.Service))),
	),

	fx.Provide(teamserver.New),
)
