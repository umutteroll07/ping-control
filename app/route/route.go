package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/umutteroll07/new-task/app/handlers"
)

func SetupRoutes(app *fiber.App){

	app.Get("/ping/:ip/:count",handlers.TestPing)

}