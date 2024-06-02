package utils

import "strings"

func RemoveCommand(text *string) {
	// Remove the "/imagem" command from the incoming text
	*text = strings.TrimSpace(strings.Replace(*text, "/imagem", "", 1))
}
