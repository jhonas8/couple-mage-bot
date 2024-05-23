package actions

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ImediatellyReplyUser(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Estou pensando 🧠...")

	msg.ReplyToMessageID = update.Message.MessageID

	bot.Send(msg)
}