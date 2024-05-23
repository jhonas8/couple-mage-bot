package utils

import (
	"strings"
)

func ExtractCommand(text string) (string, bool) {
	// extract the first word from the text
	parts := strings.Split(text, " ")
	if len(parts) == 0 || !strings.HasPrefix(parts[0], "/") {
		return "", false
	}
	return parts[0], true
}
