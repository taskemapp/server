package main

import (
	"context"
	"github.com/taskemapp/server/apps/server/internal/app"
	"go.uber.org/fx"
)

func main() {
	a := fx.New(
		app.App,
	)

	a.Run()

	defer func(app *fx.App, ctx context.Context) {
		err := app.Stop(ctx)
		panic(err)
	}(a, context.Background())

	<-a.Done()
}
