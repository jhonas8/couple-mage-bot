package commands

import (
	"couplebot/clients"
	"fmt"
)

func ShowAllMovies(msgText *string) string {
	movies := clients.GetAllMovies()

	if len(movies) <= 0 {
		msg := "Não consegui achar nenhum filme em minha base de dados.\n Talvez você não tenha mais filmes salvos."
		*msgText = msg
		return msg
	}

	allMovies := ""

	for i, m := range movies {
		allMovies += fmt.Sprintf("%d. %s\n", i+1, m.Name)
	}

	msg := "Você tem esses filmes salvos: \n" + allMovies

	*msgText = msg

	return msg
}
