package main

import "github.com/isaac/app"

func main() {
	var app app.App
	app.Initialize()
	app.InitializeRoutes()
	app.Run("9000")
}
