package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/umutteroll07/new-task/app/route"
)


func main() {
	app := *fiber.New()
	route.SetupRoutes(&app)
	app.Listen(":3010")
}
