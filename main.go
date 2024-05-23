package main

import (
	"couplebot/clients"
	"couplebot/commands"
	"couplebot/utils"
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func botExecution() {
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
				clients.ImediatellyMentionUser(bot, update)

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
				msg.ReplyToMessageID = update.Message.MessageID

				command, isCommand := utils.ExtractCommand(update.Message.Text)

				if isCommand {
					switch command {
					case "/image":
						commands.ProcessImageGeneration(update.Message.Text, bot, update.Message.Chat.ID)

					}
				}

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

func serveHealthCheck() {
	log.Print("Starting health check")

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error serving health check: %s", err)
	}
}

func main() {
	botExecution()
	serveHealthCheck()
}
