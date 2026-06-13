package account

import (
	"errors"
	"ghivarra/afterglow/src/exceptions"
	"ghivarra/afterglow/src/mapping/api"
	"ghivarra/afterglow/src/mapping/entity/dbentity"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
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
		return exceptions.NewAppException(
			400,
			"Failed to parse request body",
			err,
			payload,
		)
	}

	result, err := createAccount(payload)
	if err != nil {
		return exceptions.NewAppException(
			500,
			"Failed to create account",
			err,
			payload,
		)
	}

	return ctx.Status(201).JSON(api.Response[AccountCreateResponseDto, *string]{
		Status:  "success",
		Message: "Account created successfully",
		Data:    result,
		Errors:  nil,
	})
}

// Update partially updates an account
// @Summary      Update Account
// @Description  Partially updates account credentials. Id, username, and isActive cannot be updated here
// @Tags         Account
// @Accept       json
// @Produce      json
// @Security     ApiKey
// @Param        id       path      string                   true  "Account id"
// @Param        payload  body      AccountUpdateRequestDto  true  "Account update payload"
// @Success      200      {object}  api.Response[dbentity.AccountEntity, string]
// @Failure      400      {object}  api.Response[string, string]
// @Failure      404      {object}  api.Response[string, string]
// @Failure      500      {object}  api.Response[string, string]
// @Router       /api/account/{id} [patch]
func Update(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	var payload AccountUpdateRequestDto

	if err := ctx.Bind().Body(&payload); err != nil {
		return exceptions.NewAppException(
			400,
			"Failed to parse request body",
			err,
			payload,
		)
	}

	result, err := updateAccount(id, payload)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return exceptions.NewAppException(
				404,
				"Account is not found",
				err,
				id,
			)
		}

		return exceptions.NewAppException(
			500,
			"Failed to update account",
			err,
			payload,
		)
	}

	return ctx.Status(200).JSON(api.Response[*dbentity.AccountEntity, *string]{
		Status:  "success",
		Message: "Account updated successfully",
		Data:    result,
		Errors:  nil,
	})
}

// Delete deletes an account
// @Summary      Delete Account
// @Description  Deletes an account by id
// @Tags         Account
// @Accept       json
// @Produce      json
// @Security     ApiKey
// @Param        id   path      string  true  "Account id"
// @Success      200  {object}  api.Response[string, string]
// @Failure      404  {object}  api.Response[string, string]
// @Failure      500  {object}  api.Response[string, string]
// @Router       /api/account/{id} [delete]
func Delete(ctx fiber.Ctx) error {
	id := ctx.Params("id")

	err := deleteAccount(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return exceptions.NewAppException(
				404,
				"Account is not found",
				err,
				id,
			)
		}

		return exceptions.NewAppException(
			500,
			"Failed to delete account",
			err,
			id,
		)
	}

	return ctx.Status(200).JSON(api.Response[string, *string]{
		Status:  "success",
		Message: "Account deleted successfully",
		Data:    "OK",
		Errors:  nil,
	})
}
