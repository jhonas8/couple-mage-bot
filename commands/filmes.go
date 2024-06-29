package commands

import (
	"couplebot/clients"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ShowAllMovies(msgText *string, bot *tgbotapi.BotAPI, chatID int64) {
	movies := clients.GetAllMovies()

	if len(movies) <= 0 {
		msg := tgbotapi.NewMessage(chatID, "Não consegui achar nenhum filme em minha base de dados.\nTalvez você não tenha mais filmes salvos.")
		bot.Send(msg)
		return
	}

	for i, m := range movies {
		caption := fmt.Sprintf("%d. %s", i+1, m.Title)
		photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(m.Poster))
		photo.Caption = caption
		_, err := bot.Send(photo)
		if err != nil {
			// Handle error (e.g., log it or send a text message instead)
			textMsg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s (Poster indisponível)", caption))
			bot.Send(textMsg)
		}
	}

	summaryMsg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Total de filmes: %d", len(movies)))
	bot.Send(summaryMsg)
}

func PromptForManualEntry(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Nenhum filme encontrado. Por favor, digite o título completo do filme:")
	bot.Send(msg)
}
