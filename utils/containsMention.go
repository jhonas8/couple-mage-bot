package utils

import (
	"log"
	"strings"
)

func ContainsMention(messageText string, botUsername string) bool {
	log.Printf("\nMessage text is: %s", messageText)
	log.Printf("\nBot username is: %s", botUsername)

	return strings.Contains(messageText, "@"+botUsername)
}
