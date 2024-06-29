package utils

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SetEnvironment(chatId *int64, command *string, isKnown *bool, update *tgbotapi.Update) {
	if update.Message != nil {
		commandNew, _ := ExtractCommand(update.Message.Text)
		*command = commandNew
		*chatId = update.Message.Chat.ID
		*isKnown = true
	} else if update.CallbackQuery != nil {
		*command = update.CallbackQuery.Data
		*chatId = update.CallbackQuery.Message.Chat.ID
		*isKnown = true
	} else {
		*isKnown = false
	}
}

func SetLocalText(localText *string, update *tgbotapi.Update) {
	if update.CallbackQuery != nil {
		*localText = update.CallbackQuery.Data
	} else if update.Message != nil {
		*localText = update.Message.Text
	}

}
