package taskfx

import (
	"github.com/taskemapp/server/apps/server/internal/repository/task"
	"go.uber.org/fx"
)

var App = fx.Options(
	fx.Module(
		"task",
		fx.Provide(
			fx.Private,
			fx.Annotate(task.NewPgx, fx.As(new(task.Repository))),
		),
	),
)
