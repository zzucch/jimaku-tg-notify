package storage

import (
	"errors"

	"gorm.io/gorm"
)

type User struct {
	ChatID                   int64 `gorm:"primaryKey"`
	NotificationInterval     int
	APIKey                   string
	UTCOffsetMinutes         int
	LastUpdateCheckTimestamp int64
}

func (s *Storage) AddOrGetUser(chatID int64) (User, error) {
	var existingUser User

	err := s.db.Where("chat_id = ?", chatID).First(&existingUser).Error
	if err == nil {
		return existingUser, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return existingUser, err
	}

	user := User{
		ChatID:               chatID,
		NotificationInterval: defaultInterval,
		APIKey:               defaultAPIKey,
		UTCOffsetMinutes:     0,
	}

	return user, s.db.Create(&user).Error
}

func (s *Storage) SetAPIKey(chatID int64, apiKey string) error {
	var user User
	if err := s.db.First(
		&user,
		"chat_id = ?",
		chatID).Error; err != nil {
		return errors.New("User not found")
	}

	user.APIKey = apiKey

	if err := s.db.Save(&user).Error; err != nil {
		return errors.New("Failed to update API key")
	}

	return nil
}

func (s *Storage) GetAPIKey(chatID int64) (string, error) {
	var user User
	if err := s.db.First(
		&user,
		"chat_id = ?",
		chatID).Error; err != nil {
		return "", errors.New("User not found")
	}

	return user.APIKey, nil
}

func (s *Storage) SetNotificationInterval(chatID int64, interval int) error {
	if interval <= 0 {
		return errors.New("Notification interval must be greater than 0")
	}

	var user User
	if err := s.db.First(
		&user,
		"chat_id = ?",
		chatID).Error; err != nil {
		return errors.New("User not found")
	}

	user.NotificationInterval = interval

	if err := s.db.Save(&user).Error; err != nil {
		return errors.New("Failed to update notification interval")
	}

	return nil
}

func (s *Storage) GetAllUsers() ([]User, error) {
	var users []User

	if err := s.db.Model(&User{}).Find(&users).Error; err != nil {
		return nil, errors.New("Failed to get users")
	}

	return users, nil
}

func (s *Storage) SetLastUpdateCheck(chatID int64, timestamp int64) error {
	var user User
	if err := s.db.First(&user, "chat_id = ?", chatID).Error; err != nil {
		return errors.New("User not found")
	}

	user.LastUpdateCheckTimestamp = timestamp

	if err := s.db.Save(&user).Error; err != nil {
		return errors.New("Failed to update last update check timestamp")
	}

	return nil
}

func (s *Storage) GetLastUpdateCheck(chatID int64) (int64, error) {
	var user User
	if err := s.db.First(&user, "chat_id = ?", chatID).Error; err != nil {
		return 0, errors.New("User not found")
	}

	return user.LastUpdateCheckTimestamp, nil
}

func (s *Storage) SetUTCOffset(chatID int64, offsetMinutes int) error {
	var user User
	if err := s.db.First(&user, "chat_id = ?", chatID).Error; err != nil {
		return errors.New("User not found")
	}

	user.UTCOffsetMinutes = offsetMinutes

	if err := s.db.Save(&user).Error; err != nil {
		return errors.New("Failed to update UTC offset")
	}

	return nil
}

func (s *Storage) GetUTCOffset(chatID int64) (int, error) {
	var user User
	if err := s.db.First(&user, "chat_id = ?", chatID).Error; err != nil {
		return 0, errors.New("User not found")
	}

	return user.UTCOffsetMinutes, nil
}
