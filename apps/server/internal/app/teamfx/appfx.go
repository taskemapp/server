package teamfx

import (
	teamserver "github.com/taskemapp/server/apps/server/internal/grpc/team"
	"github.com/taskemapp/server/apps/server/internal/repository/team"
	"github.com/taskemapp/server/apps/server/internal/repository/team_member"
	"github.com/taskemapp/server/apps/server/internal/repository/token"
	teamservice "github.com/taskemapp/server/apps/server/internal/service/team"
	v1 "github.com/taskemapp/server/apps/server/tools/gen/grpc/v1"
	"go.uber.org/fx"
)

var App = fx.Options(
	fx.Module(
		"team",
		fx.Provide(
			fx.Private,
			fx.Annotate(team.NewPgx, fx.As(new(team.Repository))),
			fx.Annotate(token.NewClient, fx.As(new(token.Repository))),
			fx.Annotate(team_member.NewPgx, fx.As(new(team_member.Repository))),
			fx.Annotate(teamservice.New, fx.As(new(teamservice.Service))),
		),

		fx.Provide(
			fx.Annotate(teamserver.New, fx.As(new(v1.TeamServer))),
		),
	),
)
