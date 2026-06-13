package accounttoken

import (
	"ghivarra/afterglow/src/exceptions"
	"ghivarra/afterglow/src/mapping/api"

	"github.com/gofiber/fiber/v3"
)

// Update updates account token
// @Summary      Update Account Token
// @Description  Authenticates to Contabo using the supplied username and stores the access token
// @Tags         Account Token
// @Accept       json
// @Produce      json
// @Security     ApiKey
// @Param        payload  body      UpdateTokenRequestDto  true  "Update token payload"
// @Success      200      {object}  api.Response[string, string]
// @Failure      400      {object}  api.Response[string, string]
// @Failure      500      {object}  api.Response[string, string]
// @Router       /api/account-token/update [put]
func Update(ctx fiber.Ctx) error {
	var payload UpdateTokenRequestDto

	if err := ctx.Bind().Body(&payload); err != nil {
		return exceptions.NewAppException(
			400,
			"Failed to parse request body",
			err,
			payload,
		)
	}

	result, err := updateToken(payload)
	if err != nil {
		return exceptions.NewAppException(
			500,
			"Failed to update account token",
			err,
			payload,
		)
	}

	return ctx.Status(200).JSON(api.Response[string, *string]{
		Status:  "success",
		Message: "Account token updated successfully",
		Data:    result,
		Errors:  nil,
	})
}
