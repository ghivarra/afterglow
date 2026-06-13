package accounttable

import (
	"ghivarra/afterglow/configuration/database"
	"ghivarra/afterglow/src/mapping/entity/dbentity"
)

func FetchById(id string) (*dbentity.AccountEntity, error) {
	var result dbentity.AccountEntity
	var mainTable = dbentity.AccountTableName

	res := database.Ctx.
		Table(mainTable).
		Where("id = ?", id).
		First(&result)

	return &result, res.Error
}
