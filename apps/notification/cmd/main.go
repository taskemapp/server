package main

import (
	"context"
	"go.uber.org/fx"
	"log"
	"notification/internal/app"
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
