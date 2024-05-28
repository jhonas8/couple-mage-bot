package utils

import (
	"log"
	"strings"
)

func RemovesMention(messageText string, botUsername string) string {
	log.Printf("\nOriginal message text: %s", messageText)
	log.Printf("\nBot username: %s", botUsername)

	mention := "@" + botUsername
	if strings.Contains(messageText, mention) {
		messageText = strings.ReplaceAll(messageText, mention, "")
		log.Printf("\nUpdated message text: %s", messageText)
	}

	return messageText
}
