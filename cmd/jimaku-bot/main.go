package main

import (
	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"github.com/zzucch/jimaku-tg-notify/internal/bot"
	"github.com/zzucch/jimaku-tg-notify/internal/config"
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

	b, err := bot.Initialize(config)
	if err != nil {
		log.Fatal("failed to initialize bot", "err", err)
	}

	b.Start()
}
