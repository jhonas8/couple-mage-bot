package clients

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/vertexai/genai"
)

type GenerationResponse struct {
	Candidates []struct {
		Index   int `json:"Index"`
		Content struct {
			Role  string   `json:"Role"`
			Parts []string `json:"Parts"`
		} `json:"Content"`
		FinishReason  int `json:"FinishReason"`
		SafetyRatings []struct {
			Category         int     `json:"Category"`
			Probability      int     `json:"Probability"`
			ProbabilityScore float64 `json:"ProbabilityScore"`
			Severity         int     `json:"Severity"`
			SeverityScore    float64 `json:"SeverityScore"`
			Blocked          bool    `json:"Blocked"`
		} `json:"SafetyRatings"`
		FinishMessage    string      `json:"FinishMessage"`
		CitationMetadata interface{} `json:"CitationMetadata"`
	} `json:"Candidates"`
	PromptFeedback interface{} `json:"PromptFeedback"`
	UsageMetadata  struct {
		PromptTokenCount     int `json:"PromptTokenCount"`
		CandidatesTokenCount int `json:"CandidatesTokenCount"`
		TotalTokenCount      int `json:"TotalTokenCount"`
	} `json:"UsageMetadata"`
}

func GenerateContentFromText(targetStrings *[]string, projectID string, promptText string) error {
	location := "us-central1"
	modelName := "gemini-1.5-flash-001"

	ctx := context.Background()
	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return fmt.Errorf("error creating client: %w", err)
	}

	defer client.Close()

	gemini := client.GenerativeModel(modelName)
	prompt := genai.Text(promptText)

	resp, err := gemini.GenerateContent(ctx, prompt)
	if err != nil {
		return fmt.Errorf("error generating content: %w", err)
	}
	// See the JSON response in
	// https://pkg.go.dev/cloud.google.com/go/vertexai/genai#GenerateContentResponse.
	rb, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return fmt.Errorf("json.MarshalIndent: %w", err)
	}

	var result GenerationResponse

	jsonErrr := json.Unmarshal(rb, &result)

	if jsonErrr != nil {
		*targetStrings = []string{"Ocorreu um erro gerando a resposta"}
		return jsonErrr
	}

	for _, c := range result.Candidates {
		*targetStrings = append(*targetStrings, c.Content.Parts...)
	}

	return nil
}
