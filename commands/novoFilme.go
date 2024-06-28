package commands

import (
	"couplebot/clients"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetMovieProperties(s string) *clients.Movie {
	var m clients.Movie

	// Find the text between double quotes
	start := strings.Index(s, "\"")
	end := strings.LastIndex(s, "\"")

	if start != -1 && end != -1 && start < end {
		m.Name = s[start+1 : end]
	} else {
		// Fallback to the original logic if quotes are not found
		words := strings.Fields(s)
		if len(words) > 0 {
			if strings.HasPrefix(words[0], "/") && len(words) > 1 {
				m.Name = words[1]
			} else {
				m.Name = words[0]
			}
		}
	}

	return &m
}

func AddNewMovie(text string, msgText *string, bot *tgbotapi.BotAPI, chatID int64, OMBdMoviesAvailable []clients.OMDbMovie) {

	if len(OMBdMoviesAvailable) > 0 {
		// Create a formatted list of movies
		movieList := "Filmes encontrados:\n\n"
		for i, movie := range OMBdMoviesAvailable {
			movieList += fmt.Sprintf("%d. %s (%s)\n", i+1, movie.Title, movie.Year)
			// Add poster URL if available
			if movie.Poster != "N/A" {
				movieList += fmt.Sprintf("   %s\n", movie.Poster)
			}
			movieList += "\n"
		}

		// Create keyboard buttons
		var keyboard [][]tgbotapi.InlineKeyboardButton
		for i := range OMBdMoviesAvailable {
			keyboard = append(keyboard, []tgbotapi.InlineKeyboardButton{
				tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d", i+1), fmt.Sprintf("movie_%d", i)),
			})
		}
		keyboard = append(keyboard, []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData("Nenhum dos acima", "movie_none"),
		})

		msg := tgbotapi.NewMessage(chatID, movieList)
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(keyboard...)
		bot.Send(msg)

		*msgText = "Por favor, selecione o filme correto ou escolha 'Nenhum dos acima'."
	} else {
		PromptForManualEntry(bot, chatID)
	}
}
