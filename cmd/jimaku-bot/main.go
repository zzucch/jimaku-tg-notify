package main

import (
	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"github.com/zzucch/jimaku-tg-notify/internal/bot"
	"github.com/zzucch/jimaku-tg-notify/internal/config"
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

	bot.Start(config)
}
