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

type DeepSeekProvider struct {
	client   *http.Client
	template prompt.Template
	apiKey   string
	model    string
}

type DeepSeekRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
}

type DeepSeekResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

func NewDeepSeekProvider(apiKey, model string) *DeepSeekProvider {
	if model == "" {
		model = "deepseek-chat"
	}
	return &DeepSeekProvider{
		apiKey:   apiKey,
		model:    model,
		client:   &http.Client{},
		template: prompt.GetGlobalTemplate(),
	}
}

func (d *DeepSeekProvider) SetTemplate(template prompt.Template) {
	d.template = template
}

func (d *DeepSeekProvider) GenerateMessage(ctx context.Context, input string) (string, error) {
	if d.apiKey == "" {
		return "", fmt.Errorf("deepseek API key is required")
	}

	prompt := d.template.GeneratePrompt(input)

	request := DeepSeekRequest{
		Model: d.model,
		Messages: []Message{
			{Role: "system", Content: d.template.GetSystemPrompt()},
			{Role: "user", Content: prompt},
		},
		MaxTokens: 150,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.deepseek.com/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", d.apiKey))

	logRequest(req, body)

	resp, err := d.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var body []byte
		if b, err := io.ReadAll(resp.Body); err == nil {
			body = b
		}
		return "", fmt.Errorf("deepseek API returned status %d: %s", resp.StatusCode, string(body))
	}

	var response DeepSeekResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return response.Choices[0].Message.Content, nil
}

func (d *DeepSeekProvider) Name() string {
	return "deepseek"
}
