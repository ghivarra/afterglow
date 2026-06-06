package home

import (
	"fmt"
	"ghivarra/afterglow/configuration/database"
	"ghivarra/afterglow/environment"
	"ghivarra/afterglow/src/mapping/api"

	"github.com/gofiber/fiber/v3"
)

// Index returns the general status of the API
// @Summary      Get API Info
// @Description  Returns the environment, name, and version of the application
// @Tags         General
// @Accept       json
// @Produce      json
// @Success      200  {object}  api.Response[AppData, string]
// @Router       / [get]
func Index(ctx fiber.Ctx) error {
	return ctx.Status(200).JSON(api.Response[AppData, *string]{
		Status:  "success",
		Message: fmt.Sprintf("%s REST API is running normally", environment.APP_NAME),
		Data:    loadAppData(),
		Errors:  nil,
	})
}

// HealthCheck checks the server and database status
// @Summary      Health Check
// @Description  Checks if the API is running and can connect to the database
// @Tags         General
// @Accept       json
// @Produce      json
// @Success      200  {object}  api.Response[AppData, string]
// @Failure      500  {object}  api.Response[AppData, string]
// @Router       /health-check [get]
func HealthCheck(ctx fiber.Ctx) error {
	// check database
	var resultDB int
	tx := database.Ctx.Raw("SELECT 1").Scan(&resultDB)
	if tx.Error != nil {
		return ctx.Status(500).JSON(api.Response[AppData, string]{
			Status:  "error",
			Message: fmt.Sprintf("%s REST API is not running normally. Failed to connect to database.", environment.APP_NAME),
			Data:    loadAppData(),
			Errors:  tx.Error.Error(),
		})
	}

	return ctx.Status(200).JSON(api.Response[AppData, *string]{
		Status:  "success",
		Message: fmt.Sprintf("%s REST API is running normally and connected to the database succesfully", environment.APP_NAME),
		Data:    loadAppData(),
		Errors:  nil,
	})
}
