package taskfx

import (
	"github.com/taskemapp/server/apps/server/internal/logger"
	"github.com/taskemapp/server/apps/server/internal/repository/task"
	"go.uber.org/fx"
)

const module = "task"

var App = fx.Options(
	fx.Module(
		module,
		fx.Decorate(
			func(l logger.Logger) logger.Logger {
				return l.WithScope(module)
			},
		),
		fx.Provide(
			fx.Private,
			fx.Annotate(task.NewPgx, fx.As(new(task.Repository))),
		),
	),
)
