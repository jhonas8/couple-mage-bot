package clients

import (
	"context"
	"couplebot/utils"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

var ctx = context.Background()

func initiateGeanAi() (*genai.Client, *genai.GenerativeModel) {
	apiKey := utils.GetEnvironmentVariable("GEANAI_TOKEN", "")

	log.Printf("api key: %s", apiKey)

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))

	if err != nil {
		log.Panic(err)
	}

	newModel := client.GenerativeModel("gemini-pro")

	return client, newModel
}

func GenerateResponseForPrompt(msg *tgbotapi.MessageConfig) error {
	client, model := initiateGeanAi()

	defer client.Close()

	response, err := model.GenerateContent(ctx, genai.Text("You are a sweet assistent that tries to give the responses as precise and polite as possible. Try to not extend to much your answers as also not be way too short."), genai.Text(msg.Text))

	msg.Text = ""

	if err != nil {
		log.Fatalf("an error occurred: %s", err)
		return err
	}

	for _, cand := range response.Candidates {
		if cand.Content == nil {
			continue
		}

		for _, part := range cand.Content.Parts {

			if textPart, ok := part.(genai.Text); ok {
				msg.Text += string(textPart)
			}
		}

	}

	return nil
}
