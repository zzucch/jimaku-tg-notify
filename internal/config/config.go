package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/charmbracelet/log"
)

type Config struct {
	LogConfig
	BotConfig
}

type LogConfig struct {
	LogDebugLevel bool `env:"DEBUG_LOG,required"`
	BotDebugLevel bool `env:"BOT_DEBUG_LOG,required"`
}

type BotConfig struct {
	BotToken string `env:"BOT_TOKEN,required"`
}

func ParseEnvConfig() Config {
	config := Config{}

	if err := env.Parse(&config); err != nil {
		log.Fatal("failed to parse env config", "err", err)
	}

	return config
}
