package server

import (
	"fmt"
	"ghivarra/afterglow/configuration/routes"
	"ghivarra/afterglow/environment"
	"ghivarra/afterglow/src/handler"
	"ghivarra/afterglow/src/middleware"

	"github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"

	_ "ghivarra/afterglow/docs"
)

func CreateServer() {
	// run fiber app
	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		ErrorHandler:  handler.HandleError,
	})

	routes.CreateRouter(app)

	// serve swagger UI at /swagger/index.html
	app.Get("/swagger/*", middleware.SwaggerBasicAuth, swaggo.HandlerDefault)

	// listen
	app.Listen(fmt.Sprintf(":%d", environment.SERVER_PORT))
}
