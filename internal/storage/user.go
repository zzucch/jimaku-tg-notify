package storage

import (
	"errors"

	"gorm.io/gorm"
)

type User struct {
	ChatID               int64 `gorm:"primaryKey"`
	NotificationInterval int
	APIKey               string
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
		APIKey:               defaultApiKey,
	}

	return db.Create(&user).Error
}

func SetAPIKey(chatID int64, apiKey string) error {
	var user User
	if err := db.First(
		&user,
		"chat_id = ?",
		chatID).Error; err != nil {
		return errors.New("User not found")
	}

	user.APIKey = apiKey

	if err := db.Save(&user).Error; err != nil {
		return errors.New("Failed to update API key")
	}

	return nil
}

func GetAPIKey(chatID int64) (string, error) {
	var user User
	if err := db.First(
		&user,
		"chat_id = ?",
		chatID).Error; err != nil {
		return "", errors.New("User not found")
	}

	return user.APIKey, nil
}

func SetNotificationInterval(chatID int64, interval int) error {
	if interval <= 0 {
		return errors.New("Notification interval must be greater than 0")
	}

	var user User
	if err := db.First(
		&user,
		"chat_id = ?",
		chatID).Error; err != nil {
		return errors.New("User not found")
	}

	user.NotificationInterval = interval

	if err := db.Save(&user).Error; err != nil {
		return errors.New("Failed to update notification interval")
	}

	return nil
}

func GetAllUsers() ([]User, error) {
	var users []User

	if err := db.Model(&User{}).Find(&users).Error; err != nil {
		return nil, errors.New("Failed to get users")
	}

	return users, nil
}
