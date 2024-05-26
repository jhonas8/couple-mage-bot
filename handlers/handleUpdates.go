package handlers

import (
	"couplebot/actions"
	"couplebot/commands"
	"couplebot/utils"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func processComamnd(msg *tgbotapi.MessageConfig, update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	var msgText string

	if command, isKnown := utils.ExtractCommand(update.Message.Text); isKnown {
		switch command {
		case "/imagem":
			files := commands.ProcessImageGeneration(msg.Text, &(msg.Text), bot)

			for _, file := range files {
				sendable := tgbotapi.NewPhoto(update.Message.Chat.ID, file)
				bot.Send(sendable)
			}

		default:
			msgText = fmt.Sprintf("Ainda não sei como fazer %s. Me desculpe 🥹", command)
		}
	} else {
		msgText = fmt.Sprintf("Eu não conheço o comando %s 🤔", command)
	}

	msg.Text = msgText
}

func processDirectMentions(msg *tgbotapi.MessageConfig, bot *tgbotapi.BotAPI) {
	if commands, err := bot.GetMyCommands(); err != nil {
		msg.Text = fmt.Sprintf("Olá! Como posso te ajudar? Escolha um dos comandos existentes %s", commands)
	}

	msg.Text = "Olá! Não consigo responder a menções diretas ainda. Dê uma olhada nos meus comandos. 😄"
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
		processDirectMentions(&msg, bot)
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
