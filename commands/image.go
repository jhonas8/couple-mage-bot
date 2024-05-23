package commands

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ProcessImageGeneration(text string, bot *tgbotapi.BotAPI, chatID int64) ([]string, error) {
	newText := strings.TrimSpace(text)
	newText = strings.TrimPrefix(newText, "/")
	newText = strings.TrimSpace(newText)

	log.Printf("This is the prompt supposed to generate the image: %s", newText)

	msg := tgbotapi.NewMessage(chatID, "Ainda não consigo gerar imagens. Mas poderei em breve! :D")

	bot.Send(msg)

	return nil, nil
}
