package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"telegram-bot/models"
)

const (
	geminiEndpoint       = "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent"
	geminiRequestTimeout = 10 * time.Second
)

// geminiRequest represents the top-level request body for the Gemini API.
type geminiRequest struct {
	Contents         []geminiContent  `json:"contents"`
	GenerationConfig generationConfig `json:"generationConfig"`
}

type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type generationConfig struct {
	ResponseMimeType string       `json:"responseMimeType"`
	ResponseSchema   responseSchema `json:"responseSchema"`
}

// responseSchema tells Gemini exactly what JSON shape to return.
type responseSchema struct {
	Type       string                    `json:"type"`
	Properties map[string]schemaProperty `json:"properties"`
	Required   []string                  `json:"required"`
}

type schemaProperty struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

// geminiResponse represents the relevant fields from the Gemini API response.
type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

// geminiChallenge is the intermediate struct for unmarshaling the structured JSON
// that Gemini returns. Map this to models. Challenge after validation.
type geminiChallenge struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	CodeSnippet  string `json:"code_snippet"`
	OptionA      string `json:"option_a"`
	OptionB      string `json:"option_b"`
	OptionC      string `json:"option_c"`
	OptionD      string `json:"option_d"`
	CorrectIndex int    `json:"correct_index"`
	Hint         string `json:"hint"`
}

// buildChallengePrompt creates the prompt that instructs Gemini what to generate.
func buildChallengePrompt(difficulty string) string {
	return fmt.Sprintf(`Generate a Go programming multiple-choice challenge.
Difficulty: %s
- easy: basic Go (variables, types, loops, conditions, fmt)
- medium: slices, maps, structs, error handling, methods
- hard: goroutines, channels, interfaces, closures, concurrency

The challenge shows a Go code snippet and asks what it outputs or what's correct about it.
Provide exactly 4 answer options prefixed with "A) ", "B) ", "C) ", "D) ".
Exactly one answer must be correct.
The hint should guide without revealing the answer.`, difficulty)
}

// buildResponseSchema defines the JSON schema Gemini must follow.
// This guarantees the response matches our challenge structure.
func buildResponseSchema() responseSchema {
	return responseSchema{
		Type: "object",
		Properties: map[string]schemaProperty{
			"title":         {Type: "string", Description: "Short title for the challenge"},
			"description":   {Type: "string", Description: "The question being asked about the code"},
			"code_snippet":  {Type: "string", Description: "Go code snippet the question is about"},
			"option_a":      {Type: "string", Description: "First answer option, prefixed with A) "},
			"option_b":      {Type: "string", Description: "Second answer option, prefixed with B) "},
			"option_c":      {Type: "string", Description: "Third answer option, prefixed with C) "},
			"option_d":      {Type: "string", Description: "Fourth answer option, prefixed with D) "},
			"correct_index": {Type: "integer", Description: "Index of the correct answer (0=A, 1=B, 2=C, 3=D)"},
			"hint":          {Type: "string", Description: "A helpful hint that guides without revealing the answer"},
		},
		Required: []string{"title", "description", "code_snippet", "option_a", "option_b", "option_c", "option_d", "correct_index", "hint"},
	}
}

// GenerateChallenge calls the Gemini API to create a dynamic Go challenge.
// Returns an error if the API key is missing, the request fails, or the response is invalid.
// The caller should fall back to the static pool on error.
func GenerateChallenge(difficulty string) (*models.Challenge, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY not set")
	}

	reqBody := geminiRequest{
		Contents: []geminiContent{
			{
				Parts: []geminiPart{
					{Text: buildChallengePrompt(difficulty)},
				},
			},
		},
		GenerationConfig: generationConfig{
			ResponseMimeType: "application/json",
			ResponseSchema:   buildResponseSchema(),
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// 10-second timeout so we fall back quickly if Gemini is slow
	ctx, cancel := context.WithTimeout(context.Background(), geminiRequestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, geminiEndpoint, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-goog-api-key", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gemini request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gemini returned status %d: %s", resp.StatusCode, string(body))
	}

	var geminiResp geminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return nil, fmt.Errorf("failed to parse gemini response: %w", err)
	}

	// Extract the structured JSON text from the response
	if len(geminiResp.Candidates) == 0 ||
		len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("gemini returned empty response")
	}

	challengeJSON := geminiResp.Candidates[0].Content.Parts[0].Text

	var generated geminiChallenge
	if err := json.Unmarshal([]byte(challengeJSON), &generated); err != nil {
		return nil, fmt.Errorf("failed to parse challenge JSON: %w", err)
	}

	// Validate the generated challenge before using it
	if err := validateGeneratedChallenge(&generated); err != nil {
		return nil, fmt.Errorf("invalid generated challenge: %w", err)
	}

	challenge := &models.Challenge{
		ID:          fmt.Sprintf("llm-%d", time.Now().UnixNano()),
		Title:       generated.Title,
		Description: generated.Description,
		Difficulty:  difficulty,
		CodeSnippet: generated.CodeSnippet,
		Options: [4]string{
			generated.OptionA,
			generated.OptionB,
			generated.OptionC,
			generated.OptionD,
		},
		CorrectIdx: generated.CorrectIndex,
		Hint:       generated.Hint,
	}

	return challenge, nil
}

// validateGeneratedChallenge checks that the LLM returned sensible data.
func validateGeneratedChallenge(code *geminiChallenge) error {
	if code.CorrectIndex < 0 || code.CorrectIndex > 3 {
		return fmt.Errorf("correct_index %d out of range (must be 0-3)", code.CorrectIndex)
	}
	if code.CodeSnippet == "" {
		return fmt.Errorf("code_snippet is empty")
	}
	if code.OptionA == "" || code.OptionB == "" || code.OptionC == "" || code.OptionD == "" {
		return fmt.Errorf("one or more options are empty")
	}
	return nil
}
