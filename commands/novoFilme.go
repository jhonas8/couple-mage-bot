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

	words := strings.Split(s, " ")

	name := words[0]

	if strings.HasPrefix(name, "/") {
		name = words[1]
	}

	m.Name = name

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
