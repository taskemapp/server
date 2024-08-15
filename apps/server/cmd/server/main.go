package main

import (
	"context"
	"go.uber.org/fx"
	"server/internal/app"
	"time"
)

func main() {
	time.Sleep(time.Hour)
	a := fx.New(
		app.App,
	)

	a.Run()

	defer func(app *fx.App, ctx context.Context) {
		err := app.Stop(ctx)
		if err != nil {

		}
	}(a, context.Background())

	<-a.Done()
}
