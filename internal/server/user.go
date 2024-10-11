package server

import (
	"time"

	"github.com/charmbracelet/log"
	"github.com/zzucch/jimaku-tg-notify/internal/notification"
)

func (s *Server) AddUser(chatID int64) error {
	_, ok := s.users.LoadOrStore(chatID, struct{}{})
	if ok {
		return nil
	}

	user, err := s.store.AddOrGetUser(chatID)
	if err != nil {
		return err
	}

	s.updateCh <- notification.SchedulerUpdate{
		ChatID:   user.ChatID,
		Interval: time.Duration(user.NotificationInterval * int(time.Hour)),
	}

	return nil
}

func (s *Server) SetInterval(
	chatID int64,
	interval int,
) error {
	err := s.store.SetNotificationInterval(chatID, interval)
	if err != nil {
		log.Warn(
			"failed to set interval",
			"chatID",
			chatID,
			"err",
			err,
		)

		return err
	}

	s.updateCh <- notification.SchedulerUpdate{
		ChatID:   chatID,
		Interval: time.Duration(interval) * time.Hour,
	}

	return nil
}

func (s *Server) SetAPIKey(
	chatID int64,
	apiKey string,
) error {
	if err := s.store.SetAPIKey(chatID, apiKey); err != nil {
		log.Error(
			"failed to set api key",
			"chatID",
			chatID,
			"err",
			err,
		)

		return err
	}

	if err := s.clientManager.UpdateAPIKey(chatID); err != nil {
		log.Error(
			"failed to update api key",
			"chatID",
			chatID,
			"err",
			err,
		)

		return err
	}

	return nil
}

func (s *Server) ValidateAPIKey(chatID int64) (bool, error) {
	apiKey, err := s.store.GetAPIKey(chatID)
	if err != nil {
		log.Error(
			"failed to validate api key",
			"chatID",
			chatID,
			"err",
			err,
		)

		return false, err
	}

	return validateKey(apiKey), nil
}

func validateKey(apiKey string) bool {
	return apiKey != ""
}

func (s *Server) GetLastUpdateCheck(chatID int64) (int64, error) {
	timestamp, err := s.store.GetLastUpdateCheck(chatID)
	if err != nil {
		log.Error(
			"failed to get last update check timestamp",
			"chatID",
			chatID,
			"err",
			err,
		)
		return 0, err
	}
	return timestamp, nil
}
