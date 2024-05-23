package actions

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func DeleteMessage(msg *tgbotapi.Message, bot *tgbotapi.BotAPI, chatId int64) {
	deleteMessage := tgbotapi.NewDeleteMessage(chatId, msg.MessageID)

	bot.Send(deleteMessage)
}
