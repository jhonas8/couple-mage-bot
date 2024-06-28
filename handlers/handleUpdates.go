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

	if command, isKnown := utils.ExtractCommand(update.Message.Text); isKnown {
		log.Printf("This is the command: %s", command)
		switch command {
		case "/imagem":
			files := commands.ProcessImageGeneration(update.Message.Text, &msgText, bot)

			for _, file := range files {
				sendable := tgbotapi.NewPhoto(update.Message.Chat.ID, file)
				bot.Send(sendable)
			}

		case "/novo_filme":
			localText := update.Message.Text
			utils.RemoveCommand(&localText)
			m := commands.GetMovieProperties(localText)

			OMBdMoviesAvailable, err := clients.SearchMoviesByTitle(m.Name)

			if err != nil {
				msgText = "Ocorreu um erro ao buscar o filme: " + err.Error()
			}
			if strings.HasPrefix(update.Message.Text, "/novo_filme movie_") {
				choice := strings.TrimPrefix(update.Message.Text, "/novo_filme ")
				if choice == "movie_none" {
					commands.PromptForManualEntry(bot, update.Message.Chat.ID)
					return // Exit early as we're handling this case separately
				}

				if strings.HasPrefix(choice, "movie_") {
					index, _ := strconv.Atoi(strings.TrimPrefix(choice, "movie_"))
					if index >= 0 && index < len(OMBdMoviesAvailable) {
						selectedMovie := OMBdMoviesAvailable[index]
						err := clients.WriteNewMovie(clients.Movie{Name: selectedMovie.Title})
						if err != nil {
							msgText = "Ocorreu um erro ao adicionar o filme ao banco de dados: " + err.Error()
						} else {
							msgText = fmt.Sprintf("Filme '%s' adicionado à base de dados.", selectedMovie.Title)
						}
					}
				}
			} else {
				commands.AddNewMovie(update.Message.Text, &msgText, bot, update.Message.Chat.ID, OMBdMoviesAvailable)
			}

			sendable := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
			bot.Send(sendable)

		case "/filmes":
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
	isNotAValidInstruction := update.Message.Entities == nil || !update.Message.IsCommand() && !utils.ContainsMention(update.Message.Text, bot.Self.UserName)

	if update.Message == nil || isNotAValidInstruction {
		return
	}

	immediatelyMsg, _ := actions.ImmediatelyReplyUser(bot, *update)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.ReplyToMessageID = update.Message.MessageID

	if update.Message.IsCommand() {
		processComamnd(&msg, update, bot)
	} else if utils.ContainsMention(update.Message.Text, bot.Self.UserName) {
		processDirectMentions(&msg, update, bot)
	}

	if msg.Text == "" {
		return
	}

	utils.AddsUserMention(&msg, update.Message.From.UserName)

	actions.DeleteMessage(immediatelyMsg, bot, update.Message.Chat.ID)

	msgSent, errorSending := bot.Send(msg)

	if errorSending != nil {
		log.Printf("\nError sending message: %s", errorSending)
	} else {
		log.Printf("\nMessage sent: %s", msgSent.Text)
	}
}
