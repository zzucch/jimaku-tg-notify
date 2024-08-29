package server

import (
	"time"

	"github.com/charmbracelet/log"
	"github.com/zzucch/jimaku-tg-notify/internal/notify"
	"github.com/zzucch/jimaku-tg-notify/internal/storage"
)

func (s *Server) AddUser(chatID int64) error {
	_, ok := s.users.LoadOrStore(chatID, struct{}{})
	if ok {
		return nil
	}

	user, err := storage.AddOrGetUser(chatID)
	if err != nil {
		return err
	}

	s.updateCh <- notify.Update{
    ChatID:   user.ChatID,
		Interval: time.Duration(user.NotificationInterval * int(time.Hour)),
	}

	return nil
}

func (s *Server) SetInterval(
	chatID int64,
	interval int,
) error {
	err := storage.SetNotificationInterval(chatID, interval)
	if err != nil {
		log.Error(
			"failed to set interval",
			"chatID",
			chatID,
			"err",
			err)

		return err
	}

	return nil
}

func (s *Server) SetAPIKey(
	chatID int64,
	apiKey string,
) error {
	err := storage.SetAPIKey(chatID, apiKey)
	if err != nil {
		log.Error(
			"failed to set api key",
			"chatID",
			chatID,
			"err",
			err)

		return err
	}

	s.clientManager.UpdateAPIKey(chatID)

	return nil
}

func (s *Server) ValidateAPIKey(chatID int64) (exists bool, err error) {
	apiKey, err := storage.GetAPIKey(chatID)
	if err != nil {
		log.Error(
			"failed to validate api key",
			"chatID",
			chatID,
			"err",
			err)

		return false, err
	}

	return validateKey(apiKey), nil
}

func validateKey(apiKey string) bool {
	return apiKey != ""
}
