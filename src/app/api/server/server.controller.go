package server

import (
	"errors"
	"ghivarra/afterglow/src/exceptions"
	"ghivarra/afterglow/src/mapping/api"
	"ghivarra/afterglow/src/mapping/entity/dbentity"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// Create creates a new server
// @Summary      Create Server
// @Description  Creates a new server using the id supplied on the request body
// @Tags         Server
// @Accept       json
// @Produce      json
// @Security     ApiKey
// @Param        payload  body      ServerCreateRequestDto  true  "Server create payload"
// @Success      201      {object}  api.Response[dbentity.ServerEntity, string]
// @Failure      400      {object}  api.Response[string, string]
// @Failure      500      {object}  api.Response[string, string]
// @Router       /api/server [post]
func Create(ctx fiber.Ctx) error {
	var payload ServerCreateRequestDto

	if err := ctx.Bind().Body(&payload); err != nil {
		return exceptions.NewAppException(
			400,
			"Failed to parse request body",
			err,
			payload,
		)
	}

	result, err := createServer(payload)
	if err != nil {
		return exceptions.NewAppException(
			500,
			"Failed to create server",
			err,
			payload,
		)
	}

	return ctx.Status(201).JSON(api.Response[dbentity.ServerEntity, *string]{
		Status:  "success",
		Message: "Server created successfully",
		Data:    result,
		Errors:  nil,
	})
}

// Update partially updates a server
// @Summary      Update Server
// @Description  Partially updates server data. Id and accountId cannot be updated here
// @Tags         Server
// @Accept       json
// @Produce      json
// @Security     ApiKey
// @Param        serverId  path      int                     true  "Server id"
// @Param        payload   body      ServerUpdateRequestDto  true  "Server update payload"
// @Success      200       {object}  api.Response[dbentity.ServerEntity, string]
// @Failure      400       {object}  api.Response[string, string]
// @Failure      404       {object}  api.Response[string, string]
// @Failure      500       {object}  api.Response[string, string]
// @Router       /api/server/{serverId} [patch]
func Update(ctx fiber.Ctx) error {
	serverId := ctx.Params("serverId")
	id, err := strconv.Atoi(serverId)
	if err != nil {
		return exceptions.NewAppException(
			400,
			"Invalid server id",
			err,
			serverId,
		)
	}
	var payload ServerUpdateRequestDto

	if err := ctx.Bind().Body(&payload); err != nil {
		return exceptions.NewAppException(
			400,
			"Failed to parse request body",
			err,
			payload,
		)
	}

	result, err := updateServer(id, payload)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return exceptions.NewAppException(
				404,
				"Server is not found",
				err,
				serverId,
			)
		}

		return exceptions.NewAppException(
			500,
			"Failed to update server",
			err,
			payload,
		)
	}

	return ctx.Status(200).JSON(api.Response[*dbentity.ServerEntity, *string]{
		Status:  "success",
		Message: "Server updated successfully",
		Data:    result,
		Errors:  nil,
	})
}

// Delete deletes a server
// @Summary      Delete Server
// @Description  Deletes server data by id
// @Tags         Server
// @Accept       json
// @Produce      json
// @Security     ApiKey
// @Param        serverId  path      int  true  "Server id"
// @Success      200       {object}  api.Response[string, string]
// @Failure      404       {object}  api.Response[string, string]
// @Failure      500       {object}  api.Response[string, string]
// @Router       /api/server/{serverId} [delete]
func Delete(ctx fiber.Ctx) error {
	serverId := ctx.Params("serverId")
	id, err := strconv.Atoi(serverId)
	if err != nil {
		return exceptions.NewAppException(
			400,
			"Invalid server id",
			err,
			serverId,
		)
	}

	err = deleteServer(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return exceptions.NewAppException(
				404,
				"Server is not found",
				err,
				serverId,
			)
		}

		return exceptions.NewAppException(
			500,
			"Failed to delete server",
			err,
			serverId,
		)
	}

	return ctx.Status(200).JSON(api.Response[string, *string]{
		Status:  "success",
		Message: "Server deleted successfully",
		Data:    "OK",
		Errors:  nil,
	})
}
