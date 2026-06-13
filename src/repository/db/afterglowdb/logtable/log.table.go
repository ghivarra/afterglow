package logtable

import (
	"ghivarra/afterglow/configuration/database"
	"ghivarra/afterglow/src/mapping/entity/dbentity"
)

func CreateLog(entity dbentity.LogEntity) (dbentity.LogEntity, error) {
	res := database.Ctx.
		Table(dbentity.LogTableName).
		Create(&entity)

	return entity, res.Error
}
