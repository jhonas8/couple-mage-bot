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

	m.Name = name

	return &m
}

func AddNewMovie(text string, msgText *string, bot *tgbotapi.BotAPI) {
	utils.RemoveCommand(&text)

	m := getMovieProperties(text)

	clients.WriteNewMovie(*m)

	*msgText = fmt.Sprintf("Filme %s adicionado a base de dados", m.Name)
}
