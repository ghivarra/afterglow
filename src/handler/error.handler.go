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

	var errReason error = nil
	var errorData any = nil
	var timestamp string = ""

	if appErr, ok := errors.AsType[*exceptions.AppException](err); ok {
		code = appErr.Code
		message = appErr.Message
		errorData = appErr.ErrorData
		errReason = appErr.Reason
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
		Errors:    errReason,
		Timestamp: &timestamp,
	})
}
