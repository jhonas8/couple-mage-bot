package commands

import (
	"bytes"
	"couplebot/utils"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ResponseBody struct {
	SDGenerationJob struct {
		GenerationID  string `json:"generationId"`
		APICreditCost int    `json:"apiCreditCost"`
	} `json:"sdGenerationJob"`
}

type GeneratedImage struct {
	URL                             string        `json:"url"`
	NSFW                            bool          `json:"nsfw"`
	ID                              string        `json:"id"`
	LikeCount                       int           `json:"likeCount"`
	MotionMP4URL                    *string       `json:"motionMP4URL"`
	GeneratedImageVariationGenerics []interface{} `json:"generated_image_variation_generics"`
}

type GenerationResponse struct {
	GenerationsByPK struct {
		GeneratedImages     []GeneratedImage `json:"generated_images"`
		ModelID             string           `json:"modelId"`
		Motion              *string          `json:"motion"`
		MotionModel         *string          `json:"motionModel"`
		MotionStrength      *float64         `json:"motionStrength"`
		Prompt              string           `json:"prompt"`
		NegativePrompt      string           `json:"negativePrompt"`
		ImageHeight         int              `json:"imageHeight"`
		ImageToVideo        *string          `json:"imageToVideo"`
		ImageWidth          int              `json:"imageWidth"`
		InferenceSteps      int              `json:"inferenceSteps"`
		Seed                int              `json:"seed"`
		Public              bool             `json:"public"`
		Scheduler           string           `json:"scheduler"`
		SDVersion           string           `json:"sdVersion"`
		Status              string           `json:"status"`
		PresetStyle         string           `json:"presetStyle"`
		InitStrength        float64          `json:"initStrength"`
		GuidanceScale       float64          `json:"guidanceScale"`
		ID                  string           `json:"id"`
		CreatedAt           string           `json:"createdAt"`
		PromptMagic         bool             `json:"promptMagic"`
		PromptMagicVersion  *string          `json:"promptMagicVersion"`
		PromptMagicStrength *float64         `json:"promptMagicStrength"`
		PhotoReal           bool             `json:"photoReal"`
		PhotoRealStrength   float64          `json:"photoRealStrength"`
		FantasyAvatar       *string          `json:"fantasyAvatar"`
		GenerationElements  []interface{}    `json:"generation_elements"`
	} `json:"generations_by_pk"`
}

func extractPrompt(text *string) {
	// Remove the "/imagem" command from the incoming text
	utils.RemoveCommand(text)
}

func ProcessImageGeneration(text string, msgText *string, bot *tgbotapi.BotAPI) []tgbotapi.FileURL {
	extractPrompt(&text)
	imagesLinks, err := httpRequestMidJourney(text)

	log.Printf("Images str: %s", imagesLinks)

	if err != nil {
		log.Printf("Ocorreu um erro em gerar a imagem: %s", err)
		*msgText = fmt.Sprintf("Ocorreu um erro em gerar a imagem: %s", err)
	}

	var files []tgbotapi.FileURL
	for _, imageLink := range imagesLinks {
		file := tgbotapi.FileURL(imageLink)

		files = append(files, file)
	}

	return files
}

func httpRequestMidJourney(prompt string) ([]string, error) {
	log.Printf("This is the prompt trying to be sent: %s", prompt)

	if prompt == "" {
		log.Print("Prompt can't be empty")
		return nil, fmt.Errorf("Empty prompt")
	}

	postBody, err := json.Marshal(map[string]interface{}{
		"height":            512,
		"prompt":            prompt,
		"width":             512,
		"alchemy":           true,
		"photoReal":         true,
		"photoRealStrength": 0.5,
		"num_images":        1,
		"presetStyle":       "DYNAMIC",
	})

	if err != nil {
		log.Print("error: ", err)
		return nil, err
	}

	postBodyBuffer := bytes.NewBuffer(postBody)

	API_KEY := utils.GetEnvironmentVariable("LEONARDO_AI_KEY", "")
	BASE_URL := "https://cloud.leonardo.ai/api/rest/v1/generations"

	client := http.Client{}
	defaultHeader := http.Header{
		"content-type":  {"application/json"},
		"authorization": {"Bearer " + API_KEY},
		"accept":        {"application/json"},
	}

	request, err := http.NewRequest("POST", BASE_URL, postBodyBuffer)

	if err != nil {
		log.Print("error creating the http request: ", err)
		return nil, err
	}

	request.Header = defaultHeader.Clone()

	log.Print("making http request to generate image...")

	response, err := client.Do(request)

	if err != nil {
		log.Print("error performing the http request: ", err)
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("status code different than 200: %d", response.StatusCode)
	}

	defer response.Body.Close()

	var result ResponseBody
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		log.Print("error: ", err)

		return nil, err
	}

	var images []string

	for {
		time.Sleep(5 * time.Second)

		log.Print("making http request to check generation status...")

		checkingUrl := fmt.Sprintf("%s/%s", BASE_URL, result.SDGenerationJob.GenerationID)
		request, err := http.NewRequest("GET", checkingUrl, nil)

		if err != nil {
			log.Print("error creating the http request: ", err)
			return nil, err
		}

		request.Header = defaultHeader.Clone()

		checkingResponse, err := client.Do(request)
		if err != nil {
			log.Print("error performing the http request: ", err)
			continue
		}

		defer checkingResponse.Body.Close()

		var checkingResult GenerationResponse
		if err := json.NewDecoder(checkingResponse.Body).Decode(&checkingResult); err != nil {
			log.Print("error decoding the response: ", err)
			continue
		}

		if checkingResult.GenerationsByPK.Status == "COMPLETE" {
			log.Print("generation completed successfully")
			for _, image := range checkingResult.GenerationsByPK.GeneratedImages {
				images = append(images, image.URL)
			}
			break
		} else {
			log.Printf("generation status: %s", checkingResult.GenerationsByPK.Status)
		}
	}

	return images, nil
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
