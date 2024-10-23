package profile

import (
	profilesrv "github.com/taskemapp/server/apps/server/internal/grpc/profile"
	"github.com/taskemapp/server/apps/server/internal/repository/user_file"
	"github.com/taskemapp/server/apps/server/internal/service/profile"
	"go.uber.org/fx"
)

var App = fx.Options(
	fx.Provide(
		fx.Annotate(user_file.New, fx.As(new(user_file.Repository))),
	),

	fx.Provide(
		fx.Annotate(profile.New, fx.As(new(profile.Service))),
	),
	fx.Provide(profilesrv.New),
)
