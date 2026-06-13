package dbentity

var ServerTableName string = "server"

type ServerEntity struct {
	Id          int     `gorm:"column:id;primaryKey;" json:"id"`
	Alias       string  `gorm:"column:alias" json:"alias"`
	Name        string  `gorm:"column:name" json:"name"`
	Description *string `gorm:"column:description" json:"description"`
	IpAddress   string  `gorm:"column:ip_address" json:"ipAddress"`
	AccountId   string  `gorm:"column:account_id" json:"accountId"`
}

// set table name
func (ServerEntity) TableName() string {
	return ServerTableName
}
