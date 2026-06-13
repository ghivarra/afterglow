package handler

import (
	"errors"
	"ghivarra/afterglow/src/exceptions"
	"ghivarra/afterglow/src/mapping/api"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

// error handler
func HandleError(ctx fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"
	var errorData any = nil
	var timestamp string = ""

	var appErr *exceptions.AppException
	if errors.As(err, &appErr) {
		code = appErr.Code
		message = appErr.Message
		errorData = appErr.Errors
		timestamp = appErr.Timestamp
	} else {
		var fiberErr *fiber.Error
		timestamp = time.Now().UTC().Format(time.RFC3339)
		if errors.As(err, &fiberErr) {
			code = fiberErr.Code
			message = fiberErr.Message
		}
	}

	log.Errorf("error: %v", err)

	return ctx.Status(code).JSON(api.Response[any, any]{
		Status:    "error",
		Message:   message,
		Data:      errorData,
		Timestamp: &timestamp,
	})
}
