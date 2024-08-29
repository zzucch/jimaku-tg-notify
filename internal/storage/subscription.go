package storage

import (
	"errors"

	"gorm.io/gorm"
)

type Subscription struct {
	TitleID            int64 `gorm:"primaryKey"`
	ChatID             int64
	LatestSubtitleTime int64
}

func Subscribe(chatID, titleID, latestSubtitleTime int64) error {
	var existingSubscription Subscription

	err := db.Where(
		"chat_id = ? AND title_id = ?",
		chatID,
		titleID).First(&existingSubscription).Error
	if err == nil {
		return errors.New("Already subscribed")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("Failed to subscribe")
	}

	subscription := Subscription{
		TitleID:            titleID,
		ChatID:             chatID,
		LatestSubtitleTime: latestSubtitleTime,
	}

	if err := db.Create(&subscription).Error; err != nil {
		return errors.New("Failed to subscribe")
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
		return errors.New("No such subscription")
	} else if err != nil {
		return errors.New("Failed to unsubscribe")
	}

	if err := db.Delete(
		&Subscription{},
		"title_id = ? AND chat_id = ?",
		titleID,
		chatID).Error; err != nil {
		return errors.New("Failed to unsubscribe")
	}

	return nil
}

func GetAllSubscriptions(chatID int64) ([]Subscription, error) {
	var subscriptions []Subscription

	if err := db.Where(
		"chat_id = ?",
		chatID).Find(&subscriptions).Error; err != nil {
		return nil, errors.New("Failed to get subscriptions")
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
		return errors.New("Subscription not found")
	}

	subscription.LatestSubtitleTime = latestSubtitleTime

	if err := db.Save(&subscription).Error; err != nil {
		return errors.New("Failed to update latest subtitle timestamp")
	}

	return nil
}