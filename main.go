package main

import (
	"couplebot/clients"
	"couplebot/utils"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(utils.GetEnvironmentVariable("TELEGRAM_BOT_TOKEN", ""))

	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)

	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			botUsername := bot.Self.UserName
			if update.Message.IsCommand() || update.Message.Entities != nil && utils.ContainsMention(update.Message.Text, botUsername) {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
				msg.ReplyToMessageID = update.Message.MessageID

				clients.GenerateResponseForPrompt(&msg)

				utils.AddsUserMention(&msg, update.Message.From.UserName)

				msgSent, errorSending := bot.Send(msg)

				if errorSending != nil {
					log.Printf("\nError sending message: %s", errorSending)
				} else {
					log.Printf("\nMessage sent: %s", msgSent.Text)
				}
			} else {
				log.Printf("\nThis message is not meant to be processed by the bot")
			}
		}
	}
}
