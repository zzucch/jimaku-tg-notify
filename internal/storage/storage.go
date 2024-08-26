package storage

import (
	"errors"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const dataDir = "./_data"
const connection = "./_data/sqlite.db"

type User struct {
	ChatID int64 `gorm:"primaryKey"`
}

type Subscription struct {
	TitleID            int64 `gorm:"primaryKey"`
	ChatID             int64
	LatestSubtitleTime int64
}

var db *gorm.DB

func Start() error {
	var err error
	if err := os.MkdirAll(dataDir, os.ModePerm); err != nil {
		return err
	}

	db, err = gorm.Open(sqlite.Open(connection), &gorm.Config{})
	if err != nil {
		return err
	}

	if err := db.AutoMigrate(&User{}, &Subscription{}); err != nil {
		return err
	}

	return nil
}

func AddUser(chatID int64) error {
	user := User{ChatID: chatID}
	return db.Create(&user).Error
}

func Subscribe(chatID, titleID, latestSubtitleTime int64) error {
	subscription := Subscription{
		TitleID:            titleID,
		ChatID:             chatID,
		LatestSubtitleTime: latestSubtitleTime,
	}

	if err := db.Create(&subscription).Error; err != nil {
		return errors.New("failed to subscribe")
	}

	return nil
}

func Unsubscribe(chatID, titleID int64) error {

	if err := db.Delete(&Subscription{}, "title_id = ? AND chat_id = ?", titleID, chatID).Error; err != nil {
		return errors.New("failed to unsubscribe")
	}
	return nil

}

func GetAllChatIDs() ([]int64, error) {
	var chatIDs []int64
	if err := db.Model(&User{}).Pluck("chat_id", &chatIDs).Error; err != nil {
		return nil, errors.New("failed to get subscriptions")
	}
	return chatIDs, nil
}

func GetAllSubscriptions(chatID int64) ([]Subscription, error) {
	var subscriptions []Subscription
	if err := db.Where(
		"chat_id = ?",
		chatID).Find(&subscriptions).Error; err != nil {
		return nil, errors.New("failed to get subscriptions")
	}
	return subscriptions, nil
}
