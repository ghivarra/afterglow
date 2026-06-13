package crypto

import (
	"ghivarra/afterglow/environment"
	"ghivarra/afterglow/src/exceptions"
	"ghivarra/afterglow/src/mapping/api"

	"github.com/gofiber/fiber/v3"
)

// Decrypt decrypts text
// @Summary      Decrypt Text
// @Description  Decrypts encrypted text. This endpoint is only active in development environment
// @Tags         Crypto
// @Accept       json
// @Produce      json
// @Security     ApiKey
// @Param        payload  body      DecryptRequestDto  true  "Decrypt payload"
// @Success      200      {object}  api.Response[DecryptResponseDto, string]
// @Failure      400      {object}  api.Response[string, string]
// @Failure      403      {object}  api.Response[string, string]
// @Failure      500      {object}  api.Response[string, string]
// @Router       /api/crypto/decrypt [post]
func Decrypt(ctx fiber.Ctx) error {
	var payload DecryptRequestDto

	if err := ctx.Bind().Body(&payload); err != nil {
		return exceptions.NewAppException(
			400,
			"Failed to parse request body",
			err,
			payload,
		)
	}

	result, err := decryptText(payload)
	if err != nil {
		if environment.APP_ENV != "development" {
			return exceptions.NewAppException(
				403,
				"Forbidden",
				err,
				payload,
			)
		}

		return exceptions.NewAppException(
			500,
			"Failed to decrypt text",
			err,
			payload,
		)
	}

	return ctx.Status(200).JSON(api.Response[DecryptResponseDto, *string]{
		Status:  "success",
		Message: "Text decrypted successfully",
		Data:    result,
		Errors:  nil,
	})
}
