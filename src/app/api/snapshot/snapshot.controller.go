package snapshot

import (
	"errors"
	"ghivarra/afterglow/src/exceptions"
	"ghivarra/afterglow/src/mapping/api"
	"ghivarra/afterglow/src/mapping/entity/dbentity"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// Create creates a new snapshot
// @Summary      Create Snapshot
// @Description  Creates a new Contabo snapshot by server alias and stores the snapshot data
// @Tags         Snapshot
// @Accept       json
// @Produce      json
// @Security     ApiKey
// @Param        payload  body      SnapshotCreateRequestDto  true  "Snapshot create payload"
// @Success      201      {object}  api.Response[dbentity.SnapshotEntity, string]
// @Failure      400      {object}  api.Response[string, string]
// @Failure      404      {object}  api.Response[string, string]
// @Failure      500      {object}  api.Response[string, string]
// @Failure      503      {object}  api.Response[string, string]
// @Router       /api/snapshot [post]
func Create(ctx fiber.Ctx) error {
	var payload SnapshotCreateRequestDto

	if err := ctx.Bind().Body(&payload); err != nil {
		return exceptions.NewAppException(
			400,
			"Failed to parse request body",
			err,
			payload,
		)
	}

	result, err := createSnapshot(payload)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return exceptions.NewAppException(
				404,
				"Snapshot requirement is not found",
				err,
				payload,
			)
		}

		if errors.Is(err, errSnapshotResponseEmpty) {
			return exceptions.NewAppException(
				503,
				"Snapshot response is empty",
				err,
				payload,
			)
		}

		return exceptions.NewAppException(
			500,
			"Failed to create snapshot",
			err,
			payload,
		)
	}

	return ctx.Status(201).JSON(api.Response[dbentity.SnapshotEntity, *string]{
		Status:  "success",
		Message: "Snapshot created successfully",
		Data:    result,
		Errors:  nil,
	})
}
