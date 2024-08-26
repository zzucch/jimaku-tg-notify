package main

import (
	"time"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"github.com/zzucch/jimaku-tg-notify/internal/bot"
	"github.com/zzucch/jimaku-tg-notify/internal/config"
	"github.com/zzucch/jimaku-tg-notify/internal/notify"
	"github.com/zzucch/jimaku-tg-notify/internal/storage"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
}

func main() {
	config := config.ParseEnvConfig()
	if config.DebugLevel {
		log.SetLevel(log.DebugLevel)
	}

	log.Debug("loaded env config", "config", config)

	if err := storage.Start(); err != nil {
		log.Fatal("failed connecting to storage", "err", err)
	}

  err := bot.Initialize(config)
	if err != nil {
		log.Fatal(
			"failed initializing bot",
			"err",
			err,
			"config",
			config)
	}

	go notificationTimer()
  bot.Start()
}

func notificationTimer() {
	for {
		notify.NotifyAll()
		time.Sleep(time.Second * 10)
	}
}
