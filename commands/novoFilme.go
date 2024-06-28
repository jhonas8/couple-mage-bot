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
		// Send each movie as a separate message with an image
		for i, movie := range OMBdMoviesAvailable {
			caption := fmt.Sprintf("%d. %s (%s)", i+1, movie.Title, movie.Year)

			var msg tgbotapi.Chattable
			if movie.Poster != "N/A" {
				photoMsg := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(movie.Poster))
				photoMsg.Caption = caption
				msg = photoMsg
			} else {
				msg = tgbotapi.NewMessage(chatID, caption)
			}

			bot.Send(msg)
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

		selectionMsg := tgbotapi.NewMessage(chatID, "Por favor, selecione o filme correto:")
		selectionMsg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(keyboard...)
		bot.Send(selectionMsg)

		*msgText = "Por favor, selecione o filme correto ou escolha 'Nenhum dos acima'."
	} else {
		PromptForManualEntry(bot, chatID)
	}
}
