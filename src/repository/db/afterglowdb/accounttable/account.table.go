package accounttable

import (
	"ghivarra/afterglow/configuration/database"
	"ghivarra/afterglow/src/mapping/entity/dbentity"

	"gorm.io/gorm"
)

func CreateAccount(entity dbentity.AccountEntity) (dbentity.AccountEntity, error) {
	res := database.Ctx.
		Table(dbentity.AccountTableName).
		Create(&entity)

	return entity, res.Error
}

func FetchById(id string) (*dbentity.AccountEntity, error) {
	var result dbentity.AccountEntity
	var mainTable = dbentity.AccountTableName

	res := database.Ctx.
		Table(mainTable).
		Where("id = ?", id).
		First(&result)

	return &result, res.Error
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

func PartialUpdateAccount(id string, data map[string]any) (*dbentity.AccountEntity, error) {
	res := database.Ctx.
		Table(dbentity.AccountTableName).
		Where("id = ?", id).
		Updates(data)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return FetchById(id)
}

func DeleteAccount(id string) error {
	res := database.Ctx.
		Table(dbentity.AccountTableName).
		Where("id = ?", id).
		Delete(&dbentity.AccountEntity{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
