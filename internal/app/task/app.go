package task

import (
	"go.uber.org/fx"
	"taskem-server/internal/repositories/task"
)

var App = fx.Options(
	fx.Provide(
		fx.Annotate(task.NewPgx, fx.As(new(task.Repository))),
	),
)
