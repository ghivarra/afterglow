package servertable

import (
	"ghivarra/afterglow/configuration/database"
	"ghivarra/afterglow/src/mapping/entity/dbentity"

	"gorm.io/gorm"
)

func CreateServer(entity dbentity.ServerEntity) (dbentity.ServerEntity, error) {
	res := database.Ctx.
		Table(dbentity.ServerTableName).
		Create(&entity)

	return entity, res.Error
}

func FetchById(id int) (*dbentity.ServerEntity, error) {
	var result dbentity.ServerEntity
	var mainTable = dbentity.ServerTableName

	res := database.Ctx.
		Table(mainTable).
		Where("id = ?", id).
		First(&result)

	return &result, res.Error
}

func FetchByAlias(alias string) (*dbentity.ServerEntity, error) {
	var result dbentity.ServerEntity
	var mainTable = dbentity.ServerTableName

	res := database.Ctx.
		Table(mainTable).
		Where("alias = ?", alias).
		First(&result)

	return &result, res.Error
}

func PartialUpdateServer(id int, data map[string]any) (*dbentity.ServerEntity, error) {
	res := database.Ctx.
		Table(dbentity.ServerTableName).
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

func DeleteServer(id int) error {
	res := database.Ctx.
		Table(dbentity.ServerTableName).
		Where("id = ?", id).
		Delete(&dbentity.ServerEntity{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
