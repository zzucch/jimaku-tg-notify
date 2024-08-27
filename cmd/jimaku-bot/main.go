package main

import (
	"time"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"github.com/zzucch/jimaku-tg-notify/internal/bot"
	"github.com/zzucch/jimaku-tg-notify/internal/client"
	"github.com/zzucch/jimaku-tg-notify/internal/config"
	"github.com/zzucch/jimaku-tg-notify/internal/notify"
	"github.com/zzucch/jimaku-tg-notify/internal/server"
	"github.com/zzucch/jimaku-tg-notify/internal/storage"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
}

func main() {
	config := config.ParseEnvConfig()
	if config.LogDebugLevel {
		log.SetLevel(log.DebugLevel)
	}

	log.Info("loaded env config", "config", config)

	if err := storage.Start(); err != nil {
		log.Fatal("failed connecting to storage", "err", err)
	}

	users, err := storage.GetAllUsers()
	if err != nil {
		log.Fatal("failed getting users", "err", err)
	}

	client := client.NewClient(config.APIKey)

	chatIDs := make([]int64, 0, len(users))
	for _, user := range users {
		chatIDs = append(chatIDs, user.ChatID)
	}

	server := server.NewServer(chatIDs, client)

  notificationCh := make(chan notify.Notification, 1000)

	bot, err := bot.Initialize(config, server, notificationCh)
	if err != nil {
		log.Fatal("failed to initialize bot", "err", err)
	}

	manager := notify.NewNotifyManager(notificationCh, client)

	log.Debug(users)
	for _, user := range users {
		manager.AddScheduler(
			user.ChatID,
			time.Duration(int(time.Hour)*user.NotificationInterval))
	}

	bot.Start()
}
