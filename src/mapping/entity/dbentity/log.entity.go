package dbentity

import "time"

var LogTableName string = "log"

type LogEntity struct {
	Id           string     `gorm:"table:id" json:"id"`
	Url          string     `gorm:"table:url" json:"url"`
	Payload      string     `gorm:"table:payload" json:"payload"`
	ResponseCode int        `gorm:"table:response_code" json:"responseCode"`
	ResponseBody *string    `gorm:"table:response_body" json:"responseBody"`
	RequestedAt  time.Time  `gorm:"table:requested_at" json:"requestedAt"`
	RespondedAt  *time.Time `gorm:"table:responded_at" json:"respondedAt"`
}

// set table name
func (LogEntity) TableName() string {
	return LogTableName
}
