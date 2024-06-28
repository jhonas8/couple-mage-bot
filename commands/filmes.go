package commands

import (
	"couplebot/clients"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ShowAllMovies(msgText *string) string {
	movies := clients.GetAllMovies()

	if len(movies) <= 0 {
		msg := "Não consegui achar nenhum filme em minha base de dados.\n Talvez você não tenha mais filmes salvos."
		*msgText = msg
		return msg
	}

	allMovies := ""

	for i, m := range movies {
		allMovies += fmt.Sprintf("%d. %s\n", i+1, m.Title)
	}

	msg := "Você tem esses filmes salvos: \n" + allMovies

	*msgText = msg

	return msg
}

func PromptForManualEntry(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Nenhum filme encontrado. Por favor, digite o título completo do filme:")
	bot.Send(msg)
}
