package snapshottable

import (
	"ghivarra/afterglow/configuration/database"
	"ghivarra/afterglow/src/mapping/entity/dbentity"
	"time"
)

func CreateSnapshot(entity dbentity.SnapshotEntity) (dbentity.SnapshotEntity, error) {
	res := database.Ctx.
		Table(dbentity.SnapshotTableName).
		Create(&entity)

	return entity, res.Error
}

func FetchActiveSnapshotsByServerId(serverId int) ([]dbentity.SnapshotEntity, error) {
	var result []dbentity.SnapshotEntity
	var mainTable = dbentity.SnapshotTableName

	res := database.Ctx.
		Table(mainTable).
		Where("serverId = ?", serverId).
		Where("deleted_at IS NULL").
		Find(&result)

	return result, res.Error
}

func SoftDeleteSnapshotById(id string) error {
	now := time.Now().UTC().Format(time.RFC3339)

	res := database.Ctx.
		Table(dbentity.SnapshotTableName).
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Updates(map[string]any{
			"deleted_at": now,
		})

	return res.Error
}
