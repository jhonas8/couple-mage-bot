package commands

import (
	"couplebot/clients"
	"couplebot/utils"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func getMovieProperties(s string) *clients.Movie {
	var m clients.Movie

	// Find the text between double quotes
	start := strings.Index(s, "\"")
	end := strings.LastIndex(s, "\"")

	if start != -1 && end != -1 && start < end {
		m.Name = s[start+1 : end]
	} else {
		// Fallback to the original logic if quotes are not found
		words := strings.Fields(s)
		if len(words) > 0 {
			if strings.HasPrefix(words[0], "/") && len(words) > 1 {
				m.Name = words[1]
			} else {
				m.Name = words[0]
			}
		}
	}

	return &m
}

func AddNewMovie(text string, msgText *string, bot *tgbotapi.BotAPI) {
	localText := text

	utils.RemoveCommand(&localText)

	m := getMovieProperties(localText)

	err := clients.WriteNewMovie(*m)

	if err != nil {
		*msgText = "Occoreu um erro ao adicionar o seu filme ao banco de dados: \n" + err.Error()
		return
	}

	*msgText = fmt.Sprintf("Filme %s adicionado a base de dados", m.Name)

}
