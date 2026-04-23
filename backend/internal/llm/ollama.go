package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type OllamaClient struct {
	baseURL    string
	model      string
	httpClient *http.Client
}

func NewOllamaClient(baseURL string, model string) *OllamaClient {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	if model == "" {
		model = "gemma3"
	}

	return &OllamaClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		model:   model,
		httpClient: &http.Client{
			Timeout: 2 * time.Minute,
		},
	}
}

func (c *OllamaClient) Complete(ctx context.Context, prompt string) (string, error) {
	requestBody := ollamaGenerateRequest{
		Model:  c.model,
		Prompt: prompt,
		Stream: false,
	}
	payload, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("marshal ollama request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/generate", bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("create ollama request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("call ollama: %w", err)
	}
	defer resp.Body.Close()

	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read ollama response: %w", err)
	}

	var responseBody ollamaGenerateResponse
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		_ = json.Unmarshal(rawBody, &responseBody)
		if responseBody.Error != "" {
			return "", fmt.Errorf("ollama returned %d: %s", resp.StatusCode, responseBody.Error)
		}
		if trimmed := trimBody(rawBody); trimmed != "" {
			return "", fmt.Errorf("ollama returned %d: %s", resp.StatusCode, trimmed)
		}
		return "", fmt.Errorf("ollama returned %d", resp.StatusCode)
	}

	if err := json.Unmarshal(rawBody, &responseBody); err != nil {
		return "", fmt.Errorf("decode ollama response: %w; body: %s", err, trimBody(rawBody))
	}
	if responseBody.Response == "" {
		return "", fmt.Errorf("ollama returned empty response")
	}

	return responseBody.Response, nil
}

type ollamaGenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type ollamaGenerateResponse struct {
	Response string `json:"response"`
	Error    string `json:"error"`
}

func trimBody(body []byte) string {
	const max = 512
	if len(body) == 0 {
		return ""
	}
	if len(body) > max {
		body = body[:max]
	}
	return strings.TrimSpace(string(body))
}
