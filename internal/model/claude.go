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

type ClaudeProvider struct {
	client   *http.Client
	template prompt.Template
	apiKey   string
	model    string
}

type ClaudeRequest struct {
	Model     string    `json:"model"`
	System    string    `json:"system,omitempty"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ClaudeResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
}

func NewClaudeProvider(apiKey, model string) *ClaudeProvider {
	if model == "" {
		model = "claude-3-sonnet-20240229"
	}
	return &ClaudeProvider{
		apiKey:   apiKey,
		model:    model,
		client:   &http.Client{},
		template: prompt.GetGlobalTemplate(),
	}
}

func (c *ClaudeProvider) SetTemplate(template prompt.Template) {
	c.template = template
}

func (c *ClaudeProvider) GenerateMessage(ctx context.Context, input string) (string, error) {
	if c.apiKey == "" {
		return "", fmt.Errorf("claude API key is required")
	}

	prompt := c.template.GeneratePrompt(input)

	request := ClaudeRequest{
		Model:  c.model,
		System: c.template.GetSystemPrompt(),
		Messages: []Message{
			{Role: "user", Content: prompt},
		},
		MaxTokens: 150,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	logRequest(req, body)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var body []byte
		if b, err := io.ReadAll(resp.Body); err == nil {
			body = b
		}
		return "", fmt.Errorf("claude API returned status %d: %s", resp.StatusCode, string(body))
	}

	var response ClaudeResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(response.Content) == 0 {
		return "", fmt.Errorf("no content in response")
	}

	return response.Content[0].Text, nil
}

func (c *ClaudeProvider) Name() string {
	return "claude"
}
