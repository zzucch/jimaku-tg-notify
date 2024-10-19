package main

import (
	"runtime"
	"time"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"github.com/zzucch/jimaku-tg-notify/internal/bot"
	"github.com/zzucch/jimaku-tg-notify/internal/client"
	"github.com/zzucch/jimaku-tg-notify/internal/config"
	"github.com/zzucch/jimaku-tg-notify/internal/notification"
	"github.com/zzucch/jimaku-tg-notify/internal/server"
	"github.com/zzucch/jimaku-tg-notify/internal/storage"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Warn("failed to load .env", "err", err)
	}
}

func main() {
	config := config.ParseEnvConfig()

	if config.LogDebugLevel {
		log.SetLevel(log.DebugLevel)
	}

	log.Info("loaded env config", "config", config)

	store, err := storage.NewStorage()
	if err != nil {
		log.Fatal("failed to connect to storage", "err", err)
	}

	users, err := store.GetAllUsers()
	if err != nil {
		log.Fatal("failed to get users", "err", err)
	}

	chatIDs := make([]int64, 0, len(users))
	for _, user := range users {
		chatIDs = append(chatIDs, user.ChatID)
	}

	updateCh := make(chan notification.SchedulerUpdate, runtime.NumCPU())
	notificationCh := make(chan notification.Notification, runtime.NumCPU())

	clientManager := client.NewManager(store)
	server := server.NewServer(chatIDs, store, clientManager, updateCh)

	bot, err := bot.NewBot(config, server, store, notificationCh)
	if err != nil {
		log.Fatal("failed to initialize bot", "err", err)
	}

	notificationManager := notification.NewManager(
		clientManager,
		store,
		updateCh,
		notificationCh)

	go notificationManager.WatchForSchedulerUpdates()

	for _, user := range users {
		if err := notificationManager.AddScheduler(
			user.ChatID,
			time.Duration(int(time.Hour)*user.NotificationInterval),
		); err != nil {
			log.Fatal("failed to add scheduler", "user", user)
		}
	}

	bot.Start()
}
