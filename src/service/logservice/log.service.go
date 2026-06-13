package logservice

import (
	"ghivarra/afterglow/src/mapping/entity/dbentity"
	"ghivarra/afterglow/src/repository/db/afterglowdb/logtable"

	"github.com/gofiber/fiber/v3/log"
)

func StoreContaboLog(logEntity dbentity.LogEntity) {
	go func() {
		_, err := logtable.CreateLog(logEntity)
		// Handle or log error silently in the background
		if err != nil {
			log.Errorf("failed to store logger. Error: %v", err)
		}
	}()
}
