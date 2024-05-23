package clients

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ImediatellyMentionUser(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Let me think about it 🧠...")

	msg.ReplyToMessageID = update.Message.MessageID

	bot.Send(msg)
}