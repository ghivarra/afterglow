package account

import (
	"ghivarra/afterglow/src/mapping/api"

	"github.com/gofiber/fiber/v3"
)

// Create creates a new account
// @Summary      Create Account
// @Description  Creates a new active account and encrypts account credentials before saving
// @Tags         Account
// @Accept       json
// @Produce      json
// @Security     ApiKey
// @Param        payload  body      AccountCreateRequestDto  true  "Account create payload"
// @Success      201      {object}  api.Response[AccountCreateResponseDto, string]
// @Failure      400      {object}  api.Response[string, string]
// @Failure      500      {object}  api.Response[string, string]
// @Router       /api/account [post]
func Create(ctx fiber.Ctx) error {
	var payload AccountCreateRequestDto
	if err := ctx.Bind().Body(&payload); err != nil {
		return ctx.Status(400).JSON(api.Response[*string, string]{
			Status:  "error",
			Message: "Failed to parse request body",
			Errors:  err.Error(),
		})
	}

	result, err := createAccount(payload)
	if err != nil {
		return ctx.Status(500).JSON(api.Response[*string, string]{
			Status:  "error",
			Message: "Failed to create account",
			Errors:  err.Error(),
		})
	}

	return ctx.Status(201).JSON(api.Response[AccountCreateResponseDto, *string]{
		Status:  "success",
		Message: "Account created successfully",
		Data:    result,
		Errors:  nil,
	})
}
