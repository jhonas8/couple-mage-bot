package actions

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ImmediatelyReplyUser(bot *tgbotapi.BotAPI, update tgbotapi.Update) (*tgbotapi.Message, error) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Estou pensando 🧠...")

	msg.ReplyToMessageID = update.Message.MessageID

	msgSent, err := bot.Send(msg)

	if err != nil {
		log.Printf("Error trying to immediately reply to the user: %s", err)
		return nil, err
	}

	return &msgSent, nil
}
