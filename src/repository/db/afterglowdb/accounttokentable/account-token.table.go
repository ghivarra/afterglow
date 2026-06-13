package accounttokentable

import (
	"ghivarra/afterglow/configuration/database"
	"ghivarra/afterglow/src/mapping/entity/dbentity"

	"gorm.io/gorm/clause"
)

func UpsertAccountToken(entity dbentity.AccountTokenEntity) (dbentity.AccountTokenEntity, error) {
	res := database.Ctx.
		Table(dbentity.AccountTokenTableName).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "account_id"},
			},
			DoUpdates: clause.AssignmentColumns([]string{
				"access_token",
				"expired_at",
			}),
		}).
		Create(&entity)

	return entity, res.Error
}
