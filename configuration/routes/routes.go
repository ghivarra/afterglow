package routes

import (
	"ghivarra/afterglow/src/app/home"

	"github.com/gofiber/fiber/v3"
)

func CreateRouter(app *fiber.App) {
	// homepage
	app.Get("/", home.Index)
	app.Get("/health-check", home.HealthCheck)
}
