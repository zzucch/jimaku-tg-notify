package storage

import (
	"errors"

	"gorm.io/gorm"
)

type Subscription struct {
	TitleID      int64 `gorm:"primaryKey"`
	ChatID       int64
	LastModified int64
	Name         string
}

func (s *Storage) Subscribe(
	chatID, titleID, latestSubtitleTime int64,
	name string,
) error {
	var existingSubscription Subscription

	err := s.db.Where(
		"chat_id = ? AND title_id = ?",
		chatID,
		titleID).First(&existingSubscription).Error
	if err == nil {
		return errors.New("Already subscribed")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("Failed to subscribe")
	}

	subscription := Subscription{
		TitleID:      titleID,
		ChatID:       chatID,
		LastModified: latestSubtitleTime,
		Name:         name,
	}

	if err := s.db.Create(&subscription).Error; err != nil {
		return errors.New("Failed to subscribe")
	}

	return nil
}

func (s *Storage) Unsubscribe(chatID, titleID int64) error {
	var subscription Subscription

	err := s.db.Where(
		"title_id = ? AND chat_id = ?",
		titleID,
		chatID).First(&subscription).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("No such subscription")
	} else if err != nil {
		return errors.New("Failed to unsubscribe")
	}

	if err := s.db.Delete(
		&Subscription{},
		"title_id = ? AND chat_id = ?",
		titleID,
		chatID).Error; err != nil {
		return errors.New("Failed to unsubscribe")
	}

	return nil
}

func (s *Storage) SubscriptionExists(chatID, titleID int64) (bool, error) {
	var subscription Subscription

	err := s.db.Where(
		"chat_id = ? AND title_id = ?",
		chatID,
		titleID).First(&subscription).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else if err != nil {
		return false, errors.New("Failed to check subscription existence")
	}

	return true, nil
}

func (s *Storage) GetAllSubscriptions(chatID int64) ([]Subscription, error) {
	var subscriptions []Subscription

	if err := s.db.Where(
		"chat_id = ?",
		chatID).Find(&subscriptions).Error; err != nil {
		return nil, errors.New("Failed to get subscriptions")
	}

	return subscriptions, nil
}

func (s *Storage) SetLatestSubtitleTimestamp(
	chatID,
	titleID,
	latestSubtitleTime int64,
) error {
	var subscription Subscription

	if err := s.db.Where(
		"chat_id = ? AND title_id = ?",
		chatID,
		titleID).First(&subscription).Error; err != nil {
		return errors.New("Subscription not found")
	}

	subscription.LastModified = latestSubtitleTime

	if err := s.db.Save(&subscription).Error; err != nil {
		return errors.New("Failed to update latest subtitle timestamp")
	}

	return nil
}

func (s *Storage) SetName(
	chatID, titleID int64,
	name string,
) error {
	var subscription Subscription

	if err := s.db.Where(
		"chat_id = ? AND title_id = ?",
		chatID,
		titleID).First(&subscription).Error; err != nil {
		return errors.New("Subscription not found")
	}

	subscription.Name = name

	if err := s.db.Save(&subscription).Error; err != nil {
		return errors.New("Failed to update name")
	}

	return nil
}

func (s *Storage) GetSubscription(
	chatID, titleID int64,
) (*Subscription, error) {
	var subscription Subscription

	err := s.db.Where(
		"chat_id = ? AND title_id = ?",
		chatID,
		titleID).First(&subscription).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("Subscription not found")
	} else if err != nil {
		return nil, errors.New("Failed to get subscription")
	}

	return &subscription, nil
}
