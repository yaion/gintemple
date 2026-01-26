package main

import (
	"shop/internal/bootstrap"
)

func main() {
	app := bootstrap.NewApp()
	app.Run()
}
