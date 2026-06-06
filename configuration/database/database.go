package database

import (
	"ghivarra/afterglow/environment"

	"github.com/gofiber/fiber/v3/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var Ctx *gorm.DB

func ConnectDB() {
	var err error

	// set logger options
	var loggerOption logger.Interface
	if environment.APP_ENV == "development" {
		loggerOption = logger.Default.LogMode(logger.Info)
	} else {
		loggerOption = logger.Default.LogMode(logger.Error)
	}

	// connect
	Ctx, err = gorm.Open(sqlite.Open(environment.DB_PATH), &gorm.Config{
		Logger: loggerOption,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	// failed
	if err != nil {
		log.Panic("failed to connect database")
	}
}
