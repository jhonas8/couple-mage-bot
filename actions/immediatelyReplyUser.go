package actions

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ImmediatelyReplyUser(bot *tgbotapi.BotAPI, chatID int64, messageID int) (*tgbotapi.Message, error) {
	msg := tgbotapi.NewMessage(chatID, "Estou pensando 🧠...")

	msg.ReplyToMessageID = messageID

	msgSent, err := bot.Send(msg)

	if err != nil {
		log.Printf("Error trying to immediately reply to the user: %s", err)
		return nil, err
	}

	return &msgSent, nil
}
