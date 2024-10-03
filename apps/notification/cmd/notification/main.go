package main

import (
	"context"
	"github.com/taskemapp/server/apps/notification/internal/app"
	"go.uber.org/fx"
	"log"
)

func main() {
	a := fx.New(
		app.App,
	)

	a.Run()

	defer func(app *fx.App, ctx context.Context) {
		err := app.Stop(ctx)
		if err != nil {
			log.Print(err)
		}
	}(a, context.Background())

	<-a.Done()
}
