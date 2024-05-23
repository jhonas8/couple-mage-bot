package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ConfigBot(bot *tgbotapi.BotAPI) {
	log.Printf("Authorized on account %s", bot.Self.UserName)

	bot.Debug = true
}
