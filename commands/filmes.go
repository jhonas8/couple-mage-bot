package commands

import (
	"couplebot/clients"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ShowAllMovies(msgText *string, bot *tgbotapi.BotAPI, chatID int64) {
	movies := clients.GetAllMovies()

	if len(movies) <= 0 {
		msg := tgbotapi.NewMessage(chatID, "Não consegui achar nenhum filme em minha base de dados.\nTalvez você não tenha mais filmes salvos.")
		bot.Send(msg)
		return
	}

	var sentMessageIDs []int

	for i, m := range movies {
		caption := fmt.Sprintf("%d. %s", i+1, m.Title)
		photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(m.Poster))
		photo.Caption = caption

		// Create delete button
		deleteButton := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Deletar", fmt.Sprintf("delete:%s:%s", m.ImdbID, m.Title)),
			),
		)
		photo.ReplyMarkup = deleteButton

		message, err := bot.Send(photo)
		if err != nil {
			// Handle error (e.g., log it or send a text message instead)
			textMsg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s (Poster indisponível)", caption))
			message, _ = bot.Send(textMsg)
		}
		sentMessageIDs = append(sentMessageIDs, message.MessageID)
	}

	summaryMsg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Total de filmes: %d", len(movies)))
	message, _ := bot.Send(summaryMsg)
	sentMessageIDs = append(sentMessageIDs, message.MessageID)

	// Save the sent message IDs
	err := clients.SaveIdsForMovieMessages(chatID, sentMessageIDs, "all_movies")
	if err != nil {
		fmt.Printf("Error saving message IDs: %v\n", err)
	}
}

func HandleDeleteMovie(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	data := callback.Data

	// Extract the IMDb ID from the callback data
	imdbID := strings.Split(data, ":")[1]
	title := strings.Split(data, ":")[2]

	// Delete the movie from the database
	err := clients.DeleteMovieByIMDbID(imdbID)
	if err != nil {
		// Handle error
		bot.Send(tgbotapi.NewMessage(chatID, "Erro ao deletar o filme."))
		return
	}

	// Delete all previous messages
	savedIDs, documentIDs, err := clients.GetIdsForMovieMessages("all_movies")
	if err == nil {
		for _, msgID := range savedIDs {
			bot.Send(tgbotapi.NewDeleteMessage(chatID, msgID))
		}
		for _, docID := range documentIDs {
			clients.DeleteSavedIds(chatID, docID)
		}
	}

	// Show updated movie list
	ShowAllMovies(nil, bot, chatID)

	// Answer the callback query
	bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Filme deletado com sucesso: %s", title)))
}

func PromptForManualEntry(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Nenhum filme encontrado. Por favor, digite o título completo do filme:")
	bot.Send(msg)
}
