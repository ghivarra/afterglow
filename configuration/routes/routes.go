package routes

import (
	"ghivarra/afterglow/src/app/api/account"
	accounttoken "ghivarra/afterglow/src/app/api/account-token"
	"ghivarra/afterglow/src/app/home"
	"ghivarra/afterglow/src/middleware"

	"github.com/gofiber/fiber/v3"
)

func CreateRouter(app *fiber.App) {
	// homepage
	app.Get("/", home.Index)
	app.Get("/health-check", home.HealthCheck)

	// api
	app.Use("/api", middleware.ApiKeyAuth)
	app.Post("/api/account", account.Create)
	app.Patch("/api/account/:id", account.Update)
	app.Delete("/api/account/:id", account.Delete)
	app.Put("/api/account-token/update", accounttoken.Update)
}
