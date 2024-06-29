package handlers

import (
	"couplebot/actions"
	"couplebot/clients"
	"couplebot/commands"
	"couplebot/utils"
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func processComamnd(msg *tgbotapi.MessageConfig, update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	var msgText string

	if update == nil {
		log.Println("Error: update is nil")
		return
	}

	var command string
	var isKnown bool
	var chatID int64

	if update.Message != nil {
		command, isKnown = utils.ExtractCommand(update.Message.Text)
		chatID = update.Message.Chat.ID
	} else if update.CallbackQuery != nil {
		command = update.CallbackQuery.Data
		isKnown = true
		chatID = update.CallbackQuery.Message.Chat.ID
	}

	if isKnown {
		log.Printf("This is the command: %s", command)
		switch {
		case command == "/imagem":
			files := commands.ProcessImageGeneration(update.Message.Text, &msgText, bot)

			for _, file := range files {
				sendable := tgbotapi.NewPhoto(chatID, file)
				bot.Send(sendable)
			}

		case command == "/novo_filme" || strings.HasPrefix(command, "movie_"):
			var localText string

			if update.CallbackQuery != nil {
				localText = update.CallbackQuery.Data
			} else if update.Message != nil {
				localText = update.Message.Text
			}

			choice := localText
			if choice == "movie_none" {
				commands.PromptForManualEntry(bot, chatID)
				return
			}

			utils.RemoveCommand(&localText)
			m := commands.GetMovieProperties(localText)

			OMBdMoviesAvailable, err := clients.SearchMoviesByTitle(m.Title)

			if err != nil {
				msgText = "Ocorreu um erro ao buscar o filme: " + err.Error()
			}
			if strings.HasPrefix(choice, "movie_") {
				index, _ := strconv.Atoi(strings.TrimPrefix(choice, "movie_"))
				log.Printf("Index: %d", index)
				if index >= 0 && index < len(OMBdMoviesAvailable) {
					selectedMovie := OMBdMoviesAvailable[index]
					log.Printf("Selected movie: %+v", selectedMovie)
					err := clients.WriteNewMovie(selectedMovie)
					if err != nil {
						msgText = "Ocorreu um erro ao adicionar o filme ao banco de dados: " + err.Error()
					} else {
						msgText = fmt.Sprintf("Filme '%s' salvo com sucesso.", selectedMovie.Title)
					}
				}

			} else {
				commands.AddNewMovie(localText, &msgText, bot, chatID, OMBdMoviesAvailable)
			}

			sendable := tgbotapi.NewMessage(chatID, msgText)
			bot.Send(sendable)

		case command == "/filmes":
			commands.ShowAllMovies(&msgText)

			sendable := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
			bot.Send(sendable)

		default:
			msgText = fmt.Sprintf("Ainda não sei como fazer %s. Me desculpe 🥹", command)
		}
	} else {
		msgText = fmt.Sprintf("Eu não conheço o comando %s 🤔", command)
	}

	msg.Text = msgText
}

func processDirectMentions(msg *tgbotapi.MessageConfig, update *tgbotapi.Update, bot *tgbotapi.BotAPI) {

	projectID := "linen-shape-420522"

	prompt := utils.RemovesMention(update.Message.Text, bot.Self.UserName)

	prompt += "\n Se possível, formate sua resposta de forma adequada para ser enviada em mensagens de texto no telegram."

	messages := []string{}

	clients.GenerateContentFromText(&messages, projectID, prompt)

	for _, message := range messages {
		m := tgbotapi.NewMessage(update.Message.Chat.ID, message)
		bot.Send(m)
	}

	msg.Text = "\n\n Isso é tudo que eu sei sobre! 🤓"
}

func HandleUpdate(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if update.Message == nil && update.CallbackQuery == nil {
		return
	}

	var chatID int64
	var messageID int
	var userName string

	if update.Message != nil {
		chatID = update.Message.Chat.ID
		messageID = update.Message.MessageID
		userName = update.Message.From.UserName
	} else if update.CallbackQuery != nil {
		log.Printf("CallbackQuery: %+v", update.CallbackQuery)
		chatID = update.CallbackQuery.Message.Chat.ID
		messageID = update.CallbackQuery.Message.MessageID
		userName = update.CallbackQuery.From.UserName
		log.Printf("MessageID: %d", messageID)
		log.Printf("ChatID: %d", chatID)
	}

	immediatelyMsg, _ := actions.ImmediatelyReplyUser(bot, chatID, messageID)

	msg := tgbotapi.NewMessage(chatID, "")
	msg.ReplyToMessageID = messageID

	if update.Message != nil && update.Message.IsCommand() {
		processComamnd(&msg, update, bot)
	} else if update.CallbackQuery != nil {
		processComamnd(&msg, update, bot)
	} else if update.Message != nil && utils.ContainsMention(update.Message.Text, bot.Self.UserName) {
		processDirectMentions(&msg, update, bot)
	}

	actions.DeleteMessage(immediatelyMsg, bot, chatID)

	if msg.Text == "" {
		return
	}

	utils.AddsUserMention(&msg, userName)

	msgSent, errorSending := bot.Send(msg)

	if errorSending != nil {
		log.Printf("\nError sending message: %s", errorSending)
	} else {
		log.Printf("\nMessage sent: %s", msgSent.Text)
	}
}
