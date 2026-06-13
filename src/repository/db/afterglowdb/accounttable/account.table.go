package accounttable

import (
	"ghivarra/afterglow/configuration/database"
	"ghivarra/afterglow/src/mapping/entity/dbentity"
)

func CreateAccount(entity dbentity.AccountEntity) (dbentity.AccountEntity, error) {
	res := database.Ctx.
		Table(dbentity.AccountTableName).
		Create(&entity)

	return entity, res.Error
}

func FetchByUsername(username string) (*dbentity.AccountEntity, error) {
	var result dbentity.AccountEntity
	var mainTable = dbentity.AccountTableName

	res := database.Ctx.
		Table(mainTable).
		Where("username = ?", username).
		First(&result)

	return &result, res.Error
}
