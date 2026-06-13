package dbentity

import "time"

var AccountTokenTableName string = "account_token"

type AccountTokenEntity struct {
	AccountId   string    `gorm:"column:account_id;primaryKey;" json:"accountId"`
	AccessToken string    `gorm:"column:access_token" json:"accessToken"`
	ExpiredAt   time.Time `gorm:"column:expired_at" json:"expiredAt"`
}

// set table name
func (AccountTokenEntity) TableName() string {
	return AccountTokenTableName
}
