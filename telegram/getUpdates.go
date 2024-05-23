package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func GetUpdatesChannel(bot *tgbotapi.BotAPI) tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)

	u.Timeout = 120

	updates := bot.GetUpdatesChan(u)

	return updates
}
