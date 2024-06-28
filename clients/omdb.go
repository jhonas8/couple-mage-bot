package clients

import (
	"couplebot/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

type MovieSearch struct {
	Search       []OMDbMovie `json:"Search"`
	TotalResults string      `json:"totalResults"`
	Response     string      `json:"Response"`
}

// Movie represents an individual movie or series in the search results
type OMDbMovie struct {
	Title  string `json:"Title"`
	Year   string `json:"Year"`
	ImdbID string `json:"imdbID"`
	Type   string `json:"Type"`
	Poster string `json:"Poster"`
}

func getRequestUrl() string {
	apiKey := utils.GetEnvironmentVariable("OMDB_KEY", "")
	url := fmt.Sprintf("http://www.omdbapi.com/?apikey=%s", apiKey)
	return url
}

func SearchMoviesByTitle(title string) ([]OMDbMovie, error) {
	url := getRequestUrl()
	url = fmt.Sprintf("%s&s=%s", url, title)
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var movieSearch MovieSearch
	err = json.NewDecoder(resp.Body).Decode(&movieSearch)
	if err != nil {
		return nil, err
	}

	return movieSearch.Search, nil
}
