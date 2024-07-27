package main

import (
	"github.com/charmbracelet/log"
	"github.com/zzucch/jimaku-tg-notify/internal/bot"
	"github.com/zzucch/jimaku-tg-notify/internal/config"
	"github.com/joho/godotenv"
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

	log.Debug("starting bot")
	bot.Start(config)
}
