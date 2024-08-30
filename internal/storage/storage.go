package storage

import (
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	dataDir    = "./_data"
	connection = "./_data/sqlite.db"
)

const (
	defaultInterval = 6
	defaultAPIKey   = ""
)

var db *gorm.DB

func Start() error {
	var err error

	if err := os.MkdirAll(dataDir, os.ModePerm); err != nil {
		return err
	}

	db, err = gorm.Open(sqlite.Open(connection), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return err
	}

	if err := db.AutoMigrate(&User{}, &Subscription{}); err != nil {
		return err
	}

	return nil
}
