package dbentity

var AccountTokenTableName string = "account_token"

type AccountTokenEntity struct {
	AccountId   string `gorm:"column:account_id;primaryKey;" json:"accountId"`
	AccessToken string `gorm:"column:access_token" json:"accessToken"`
	ExpiredAt   string `gorm:"column:expired_at" json:"expiredAt"`
}

// set table name
func (AccountTokenEntity) TableName() string {
	return AccountTokenTableName
}
