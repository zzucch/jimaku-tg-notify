package storage

import (
	"errors"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const dataDir = "./_data"
const connection = "./_data/sqlite.db"

const defaultInterval = 6

type User struct {
	ChatID               int64 `gorm:"primaryKey"`
	NotificationInterval int   `gorm:"primaryKey"`
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

func AddUser(chatID int64) error {
	var existingUser User

	err := db.Where("chat_id = ?", chatID).First(&existingUser).Error
	if err == nil {
		return nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	user := User{
		ChatID:               chatID,
		NotificationInterval: defaultInterval,
	}

	return db.Create(&user).Error
}

func SetNotificationInterval(chatID int64, interval int) error {
	if interval <= 0 {
		return errors.New("notification interval must be greater than 0")
	}

	var user User
	if err := db.First(
		&user,
		"chat_id = ?",
		chatID).Error; err != nil {
		return errors.New("user not found")
	}

	user.NotificationInterval = interval

	if err := db.Save(&user).Error; err != nil {
		return errors.New("failed to update notification interval")
	}

	return nil
}

func Subscribe(chatID, titleID, latestSubtitleTime int64) error {
	var existingSubscription Subscription

	err := db.Where(
		"chat_id = ? AND title_id = ?",
		chatID,
		titleID).First(&existingSubscription).Error
	if err == nil {
		return errors.New("already subscribed")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("failed to subscribe")
	}

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
	var subscription Subscription

	err := db.Where(
		"title_id = ? AND chat_id = ?",
		titleID,
		chatID).First(&subscription).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("no such subscription")
	} else if err != nil {
		return errors.New("failed to unsubscribe")
	}

	if err := db.Delete(
		&Subscription{},
		"title_id = ? AND chat_id = ?",
		titleID,
		chatID).Error; err != nil {
		return errors.New("failed to unsubscribe")
	}

	return nil
}

func GetAllChatIDs() ([]int64, error) {
	var chatIDs []int64

	if err := db.Model(&User{}).Pluck(
		"chat_id",
		&chatIDs).Error; err != nil {
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

func SetLatestSubtitleTimestamp(
	chatID,
	titleID,
	latestSubtitleTime int64,
) error {
	var subscription Subscription

	if err := db.Where(
		"chat_id = ? AND title_id = ?",
		chatID,
		titleID).First(&subscription).Error; err != nil {
		return errors.New("subscription not found")
	}

	subscription.LatestSubtitleTime = latestSubtitleTime

	if err := db.Save(&subscription).Error; err != nil {
		return errors.New("failed to update latest subtitle timestamp")
	}

	return nil
}
