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

type Storage struct {
	db *gorm.DB
}

func Start() (*Storage, error) {
	var err error

	if err := os.MkdirAll(dataDir, os.ModePerm); err != nil {
		return nil, err
	}

	db, err := gorm.Open(sqlite.Open(connection), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&User{}, &Subscription{}); err != nil {
		return nil, err
	}

	return &Storage{
		db: db,
	}, nil
}
