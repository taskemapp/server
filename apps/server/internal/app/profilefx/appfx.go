package profilefx

import (
	profilesrv "github.com/taskemapp/server/apps/server/internal/grpc/profile"
	"github.com/taskemapp/server/apps/server/internal/pkg/s3"
	"github.com/taskemapp/server/apps/server/internal/repository/user"
	"github.com/taskemapp/server/apps/server/internal/repository/user_file"
	"github.com/taskemapp/server/apps/server/internal/service/profile"
	"github.com/taskemapp/server/apps/server/internal/service/profile/image"
	v1 "github.com/taskemapp/server/apps/server/tools/gen/grpc/v1"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var App = fx.Options(
	fx.Module(
		"profile",
		fx.Decorate(
			func(l *zap.Logger) *zap.Logger {
				return l.With(zap.String("scope", "profile"))
			},
		),

		fx.Provide(
			fx.Private,
			s3.NewConfig,
			s3.New,
		),

		fx.Provide(
			fx.Private,
			fx.Annotate(user.NewPgx, fx.As(new(user.Repository))),
			fx.Annotate(user_file.New, fx.As(new(user_file.Repository))),
			fx.Annotate(image.NewProcessing, fx.As(new(image.Processing))),
		),

		fx.Provide(
			fx.Private,
			fx.Annotate(profile.New, fx.As(new(profile.Service))),
		),

		fx.Invoke(s3.Invoke),
		fx.Provide(
			fx.Annotate(profilesrv.New, fx.As(new(v1.ProfileServer))),
		),
	),
)
