package dbentity

var AccountTableName string = "account"

type AccountEntity struct {
	Id           string `gorm:"column:id;primaryKey;" json:"id"`
	Username     string `gorm:"column:username" json:"username"`
	Password     string `gorm:"column:password" json:"password"`
	ApiClientId  string `gorm:"column:api_client_id" json:"apiClientId"`
	ApiClientKey string `gorm:"column:api_client_key" json:"apiClientKey"`
	IsActive     int    `gorm:"column:is_active" json:"isActive"` // 1 or 0
}

// set table name
func (AccountEntity) TableName() string {
	return AccountTableName
}
