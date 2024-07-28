package bot

import (
	"github.com/charmbracelet/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zzucch/jimaku-tg-notify/internal/config"
)

func handleMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	log.Debug("handling message", "chatID", chatID)

	sendMessage(bot, update, "idk")
}

func sendMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update, text string) {
	message := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	if _, err := bot.Send(message); err != nil {
		log.Fatal(err)
	}

	chatID := update.Message.Chat.ID
	log.Debug("sending message", "chatID", chatID, "messageText", text)
}

func Start(config config.Config) {
	log.Debug("starting bot")
	bot, err := tgbotapi.NewBotAPI(config.BotToken)
	if err != nil {
		log.Fatal(
			"failed creating BotAPI instance",
			"err",
			err,
			"token",
			config.BotConfig.BotToken)
	}

	if config.BotDebugLevel {
		bot.Debug = true
	}

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		handleMessage(bot, update)
	}
}
