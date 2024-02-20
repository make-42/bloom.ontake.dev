package main

import (
	"bloom/config"
	"bloom/db"
	"bloom/middlewares"
	"bloom/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(config.FiberConfig)

	middlewares.Init(app)
	routes.Init(app)

	db.Init()
	app.Listen(config.ListenAddress)
}
