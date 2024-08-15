package team

import (
	teamserver "github.com/taskemapp/server/apps/server/internal/grpc/team"
	"github.com/taskemapp/server/apps/server/internal/repositories/team"
	"github.com/taskemapp/server/apps/server/internal/repositories/team_member"
	teamservice "github.com/taskemapp/server/apps/server/internal/service/team"
	"go.uber.org/fx"
)

var App = fx.Options(
	fx.Provide(
		fx.Annotate(team.NewPgx, fx.As(new(team.Repository))),
		fx.Annotate(team_member.NewPgx, fx.As(new(team_member.Repository))),
		fx.Annotate(teamservice.New, fx.As(new(teamservice.Service))),
	),

	fx.Provide(teamserver.New),
)
