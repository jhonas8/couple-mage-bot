package utils

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func AddsUserMention(msg *tgbotapi.MessageConfig, targetUserMention string) {
	msg.Text = fmt.Sprintf("@%s, \n %s", targetUserMention, msg.Text)
}
