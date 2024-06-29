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
	"sync"

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

	utils.SetEnvironment(&chatID, &command, &isKnown, update)

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

			utils.SetLocalText(&localText, update)

			if localText == "movie_none" {
				commands.PromptForManualEntry(bot, chatID)
				return
			}

			utils.RemoveCommand(&localText)
			m := commands.GetMovieProperties(localText)

			OMBdMoviesAvailable, err := clients.SearchMoviesByTitle(m.Title)

			if err != nil {
				msgText = "Ocorreu um erro ao buscar o filme: " + err.Error()
			}

			if strings.HasPrefix(localText, "movie_") {
				index_and_title := strings.TrimPrefix(localText, "movie_")
				index_and_title_parts := strings.Split(index_and_title, "_")
				index, _ := strconv.Atoi(index_and_title_parts[0])
				log.Printf("Index: %d", index)
				log.Printf("Length: %d", len(OMBdMoviesAvailable))
				if index >= 0 && index < len(OMBdMoviesAvailable) {
					selectedMovie := OMBdMoviesAvailable[index]
					log.Printf("Selected movie: %+v", selectedMovie)
					err := clients.WriteNewMovie(selectedMovie)
					if err != nil {
						log.Printf("Error writing new movie: %s", err)
						msgText = "Ocorreu um erro ao adicionar o filme ao banco de dados: " + err.Error()
					} else {
						msgText = fmt.Sprintf("Filme '%s' salvo com sucesso.", selectedMovie.Title)
						sendable := tgbotapi.NewMessage(chatID, msgText)
						bot.Send(sendable)

						if selectedMovie.Poster != "N/A" {
							photoMsg := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(selectedMovie.Poster))
							bot.Send(photoMsg)
						}
					}
				}

			} else {
				commands.AddNewMovie(localText, &msgText, bot, chatID, OMBdMoviesAvailable)
			}

			log.Printf("Msg text: %s", msgText)
			if strings.Contains(msgText, "salvo com sucesso") {

				previousMessageIDs, documentIds, _ := clients.GetIdsForMovieMessages(m.Title)
				var wg sync.WaitGroup
				for _, id := range previousMessageIDs {
					wg.Add(1)
					go func(id int) {
						defer wg.Done()
						deleteMsg := tgbotapi.NewDeleteMessage(chatID, id)
						_, err := bot.Request(deleteMsg)
						if err != nil {
							log.Printf("Error deleting message %d: %v", id, err)
						}
					}(id)
				}
				wg.Wait()

				var wg2 sync.WaitGroup
				for _, id := range documentIds {
					wg2.Add(1)
					go func(id string) {
						defer wg2.Done()
						clients.DeleteSavedIds(chatID, id)
					}(id)
				}
				wg2.Wait()
			}

		case command == "/filmes":
			commands.ShowAllMovies(&msgText, bot, chatID)

			sendable := tgbotapi.NewMessage(chatID, "Estes são os filmes que você tem salvo")
			bot.Send(sendable)

		case strings.HasPrefix(command, "delete:"):
			commands.HandleDeleteMovie(bot, update.CallbackQuery)

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

	msg := tgbotapi.NewMessage(chatID, "")
	msg.ReplyToMessageID = messageID

	if update.Message != nil && update.Message.IsCommand() {
		processComamnd(&msg, update, bot)
	} else if update.CallbackQuery != nil {
		processComamnd(&msg, update, bot)
	} else if update.Message != nil && utils.ContainsMention(update.Message.Text, bot.Self.UserName) {
		processDirectMentions(&msg, update, bot)
	}

	if msg.Text == "" {
		return
	}

	immediatelyMsg, _ := actions.ImmediatelyReplyUser(bot, chatID, messageID)

	actions.DeleteMessage(immediatelyMsg, bot, chatID)

	utils.AddsUserMention(&msg, userName)

	msgSent, errorSending := bot.Send(msg)

	if errorSending != nil {
		log.Printf("\nError sending message: %s", errorSending)
	} else {
		log.Printf("\nMessage sent: %s", msgSent.Text)
	}
}
