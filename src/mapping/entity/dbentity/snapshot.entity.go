package dbentity

import "time"

var SnapshotTableName string = "snapshot"

type SnapshotEntity struct {
	Id          string     `gorm:"column:id;primaryKey;" json:"id"`
	ServerId    int        `gorm:"column:serverId" json:"serverId"`
	Name        string     `gorm:"column:name" json:"name"`
	Description *string    `gorm:"column:description" json:"description"`
	CreatedAt   time.Time  `gorm:"column:created_at" json:"createdAt"`
	DeletedAt   *time.Time `gorm:"column:deleted_at" json:"deletedAt"`
}

// set table name
func (SnapshotEntity) TableName() string {
	return SnapshotTableName
}
