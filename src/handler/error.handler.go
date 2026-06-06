package handler

import (
	"ghivarra/afterglow/src/mapping/api"

	"github.com/gofiber/fiber/v3"
)

// error handler
func HandleError(ctx fiber.Ctx, err error) error {
	// check the code
	if errData, ok := err.(*fiber.Error); ok {
		// return
		return ctx.Status(errData.Code).JSON(api.Response[*string, *string]{
			Status:  "error",
			Message: errData.Message,
		})
	}

	// return
	return ctx.Status(500).JSON(api.Response[*string, *string]{
		Status:  "error",
		Message: err.Error(),
	})
}
