package team

import (
	"go.uber.org/fx"
	teamserver "server/internal/grpc/team"
	"server/internal/repositories/team"
	"server/internal/repositories/team_member"
	teamservice "server/internal/service/team"
)

var App = fx.Options(
	fx.Provide(
		fx.Annotate(team.NewPgx, fx.As(new(team.Repository))),
		fx.Annotate(team_member.NewPgx, fx.As(new(team_member.Repository))),
		fx.Annotate(teamservice.New, fx.As(new(teamservice.Service))),
	),

	fx.Provide(teamserver.New),
)
