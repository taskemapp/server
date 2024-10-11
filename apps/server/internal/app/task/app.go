package task

import (
	"github.com/taskemapp/server/apps/server/internal/repository/task"
	"go.uber.org/fx"
)

var App = fx.Options(
	fx.Provide(
		fx.Annotate(task.NewPgx, fx.As(new(task.Repository))),
	),
)
