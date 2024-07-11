package main

import (
	"go.uber.org/fx"
	"taskem/internal/app"
)

func main() {
	fx.New(
		app.App,
	).Run()
}
