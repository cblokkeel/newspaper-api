package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type AIUtil struct {
	http   *http.Client
	apiKey string
}

func NewAIUtil(http *http.Client, apiKey string) *AIUtil {
	return &AIUtil{
		http,
		apiKey,
	}
}

// EmbeddingRequest holds the data to be sent to OpenAI
type EmbeddingRequest struct {
	Input string `json:"input"`
	Model string `json:"model"`
}

// EmbeddingResponse represents the response structure from OpenAI
type EmbeddingResponse struct {
	Data []struct {
		Object string    `json:"object"`
		Index  int       `json:"index"`
		Vector []float64 `json:"vector"`
	} `json:"data"`
}

func (ai *AIUtil) EmbedText(text string) ([]float32, error) {
	requestBody, err := json.Marshal(EmbeddingRequest{
		Input: text,
		Model: "text-embedding-ada-002",
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/embeddings", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ai.apiKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := ai.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var embeddingResponse EmbeddingResponse
	err = json.Unmarshal(body, &embeddingResponse)
	if err != nil {
		return nil, err
	}

	if len(embeddingResponse.Data) > 0 {
		vector64 := embeddingResponse.Data[0].Vector
		vector32 := make([]float32, len(vector64))
		for i, v := range vector64 {
			vector32[i] = float32(v) // Convert float64 to float32
		}
		return vector32, nil
	}

	return nil, fmt.Errorf("no embeddings returned")
}
