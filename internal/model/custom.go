package model

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/aicommit/aicommit/pkg/prompt"
)

type CustomProvider struct {
	client   *http.Client
	template prompt.Template
	apiKey   string
	model    string
	url      string
}

func NewCustomProvider(url, apiKey, model string) *CustomProvider {
	return &CustomProvider{
		apiKey:   apiKey,
		model:    model,
		url:      url,
		client:   &http.Client{},
		template: prompt.GetGlobalTemplate(),
	}
}

func (c *CustomProvider) GenerateCommitMessage(ctx context.Context, diff string) (string, error) {
	if c.url == "" {
		return "", fmt.Errorf("custom provider URL is required")
	}

	promptStr := c.template.GeneratePrompt(diff)

	// Use standard OpenAI chat format as it's the most common
	request := OpenAIRequest{
		Model: c.model,
		Messages: []Message{
			{Role: "system", Content: c.template.GetSystemPrompt()},
			{Role: "user", Content: promptStr},
		},
	}

	body, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to %s: %w", c.url, err)
	}
	defer resp.Body.Close()

	var responseBody []byte
	if resp.StatusCode != http.StatusOK {
		if body, err := io.ReadAll(resp.Body); err == nil {
			responseBody = body
		}
		return "", fmt.Errorf("custom API returned status %d: %s", resp.StatusCode, string(responseBody))
	}

	// Read response body
	responseBody, err = io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var response ChatCompletionResponse
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w, body: %s", err, string(responseBody))
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no choices in response from custom provider")
	}

	choice := response.Choices[0]
	content := choice.Message.Content

	if content == "" {
		return "", fmt.Errorf("custom provider returned empty content")
	}

	return content, nil
}

func (c *CustomProvider) Name() string {
	return "custom"
}
