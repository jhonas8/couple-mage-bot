package main

import (
	"couplebot/handlers"
	"couplebot/telegram"
	"couplebot/utils"
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func botExecution() {
	bot, err := tgbotapi.NewBotAPI(utils.GetEnvironmentVariable("TELEGRAM_BOT_TOKEN", ""))
	if err != nil {
		log.Panic("Error Creating the bot: ", err)
	}

	telegram.ConfigBot(bot)

	updates := telegram.GetUpdatesChannel(bot)

	for update := range updates {
		handlers.HandleUpdate(&update, bot)
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
	go serveHealthCheck()
	botExecution()

}
