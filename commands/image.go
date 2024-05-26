package commands

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ResponseBody struct {
	Images []string `json:"images"`
}

func ProcessImageGeneration(text string, msgText *string, bot *tgbotapi.BotAPI) []tgbotapi.FileBytes {

	imagesBase64, err := httpRequestMidJourney(text)

	if err != nil {
		*msgText = fmt.Sprintf("Ocorreu um erro em gerar a imagem: %s", err)
	}

	imagesBytes, err := decodeBase64Strings(imagesBase64)

	if err != nil {
		log.Fatal("error: ", err)
	}

	var files []tgbotapi.FileBytes

	for i, fbytes := range imagesBytes {
		file := tgbotapi.FileBytes{
			Name:  fmt.Sprintf("image_%d.jpg", i),
			Bytes: fbytes,
		}

		files = append(files, file)
	}

	return files
}

func httpRequestMidJourney(prompt string) ([]string, error) {
	postBody, err := json.Marshal(map[string]string{
		"prompt": prompt,
	})

	if err != nil {
		return nil, err
	}

	postBodyBuffer := bytes.NewBuffer(postBody)

	response, err := http.Post("https://midjourney-scrapper-f2upl4bura-uc.a.run.app/generate", "application/json", postBodyBuffer)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var result ResponseBody
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Images, nil
}

func decodeBase64Strings(imgsBase64 []string) ([][]byte, error) {
	var imgsBytes [][]byte

	for _, imgString := range imgsBase64 {
		data, err := base64.StdEncoding.DecodeString(imgString)
		if err != nil {
			return nil, fmt.Errorf("failed to decode base64 string: %w", err)
		}
		imgsBytes = append(imgsBytes, data)
	}

	return imgsBytes, nil
}
