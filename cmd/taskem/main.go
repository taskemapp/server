package main

import (
	"context"
	"go.uber.org/fx"
	"taskem/internal/app"
)

func main() {
	app := fx.New(
		app.App,
	)

	app.Run()

	defer app.Stop(context.Background())

	app.Run()

	<-app.Done()
}
