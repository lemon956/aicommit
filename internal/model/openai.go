package model

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/aicommit/aicommit/pkg/prompt"
)

type OpenAIProvider struct {
	apiKey   string
	model    string
	client   *http.Client
	template prompt.Template
}

type OpenAIRequest struct {
	Model               string    `json:"model"`
	Messages            []Message `json:"messages"`
	MaxTokens           int       `json:"max_tokens,omitempty"`
	MaxCompletionTokens int       `json:"max_completion_tokens,omitempty"`
}

// OpenAIListModelsResponse represents the response from OpenAI models list API
type OpenAIListModelsResponse struct {
	Object  string            `json:"object"`
	Data    []OpenAIModelInfo `json:"data"`
	FirstID string            `json:"first_id"`
	LastID  string            `json:"last_id"`
	HasMore bool              `json:"has_more"`
}

// OpenAIModelInfo represents individual model information
type OpenAIModelInfo struct {
	Object            string      `json:"object"`
	ID                string      `json:"id"`
	Model             string      `json:"model"`
	Created           int         `json:"created"`
	RequestID         string      `json:"request_id"`
	ToolChoice        interface{} `json:"tool_choice"`
	Usage             TokenUsage  `json:"usage"`
	Seed              int64       `json:"seed"`
	TopP              float64     `json:"top_p"`
	Temperature       float64     `json:"temperature"`
	PresencePenalty   float64     `json:"presence_penalty"`
	FrequencyPenalty  float64     `json:"frequency_penalty"`
	SystemFingerprint string      `json:"system_fingerprint"`
	InputUser         interface{} `json:"input_user"`
	ServiceTier       string      `json:"service_tier"`
	Tools             interface{} `json:"tools"`
	Metadata          interface{} `json:"metadata"`
	Choices           []Choice    `json:"choices"`
	ResponseFormat    interface{} `json:"response_format"`
}

// TokenUsage represents token usage information
type TokenUsage struct {
	TotalTokens      int `json:"total_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	PromptTokens     int `json:"prompt_tokens"`
}

// Choice represents a completion choice
type Choice struct {
	Index        int         `json:"index"`
	Message      Message     `json:"message"`
	FinishReason string      `json:"finish_reason"`
	Logprobs     interface{} `json:"logprobs"`
}

// ChatCompletionResponse represents the chat completion API response
type ChatCompletionResponse struct {
	ID                string     `json:"id"`
	Object            string     `json:"object"`
	Created           int        `json:"created"`
	Model             string     `json:"model"`
	Choices           []Choice   `json:"choices"`
	Usage             TokenUsage `json:"usage"`
	ServiceTier       string     `json:"service_tier"`
	SystemFingerprint string     `json:"system_fingerprint"`
}

func NewOpenAIProvider(apiKey, model string) *OpenAIProvider {
	if model == "" {
		model = "gpt-3.5-turbo"
	}
	return &OpenAIProvider{
		apiKey:   apiKey,
		model:    model,
		client:   &http.Client{},
		template: prompt.GetGlobalTemplate(),
	}
}

func (o *OpenAIProvider) SetTemplate(template prompt.Template) {
	o.template = template
}

// OpenAIModelsResponse represents the response from OpenAI models list API
type OpenAIModelsResponse struct {
	Object string `json:"object"`
	Data   []struct {
		ID      string `json:"id"`
		Object  string `json:"object"`
		Created int    `json:"created"`
		OwnedBy string `json:"owned_by"`
	} `json:"data"`
}

// modelCache caches available models to reduce API calls
var modelCache = struct {
	models    []string
	cacheTime time.Time
	mu        sync.Mutex
}{}

const cacheDuration = 24 * time.Hour

// validateOpenAIModel validates the model name by checking against OpenAI's models API
func (o *OpenAIProvider) validateOpenAIModel(model string) error {
	if model == "" {
		return fmt.Errorf("model name cannot be empty")
	}

	if strings.ContainsAny(model, " \t\n\r") {
		return fmt.Errorf("model name cannot contain whitespace")
	}

	// 对于特定的自定义模型，直接允许通过
	if model == "gpt-5-nano-2025-08-07" {
		return nil
	}

	// 首先检查缓存
	modelCache.mu.Lock()
	if time.Since(modelCache.cacheTime) < cacheDuration {
		for _, cachedModel := range modelCache.models {
			if cachedModel == model {
				modelCache.mu.Unlock()
				return nil
			}
		}
	}
	modelCache.mu.Unlock()

	// 尝试从OpenAI API获取模型列表
	if err := o.checkModelExists(model); err != nil {
		// 如果API调用失败，回退到基本的格式验证
		// 允许用户自定义模型，如 gpt-5-nano-2025-08-07
		return nil
	}

	return nil
}

// checkModelExists checks if a model exists by calling OpenAI's models API
func (o *OpenAIProvider) checkModelExists(model string) error {
	// 首先检查缓存
	modelCache.mu.Lock()
	if time.Since(modelCache.cacheTime) < cacheDuration {
		for _, cachedModel := range modelCache.models {
			if cachedModel == model {
				modelCache.mu.Unlock()
				return nil
			}
		}
		// 如果缓存存在但模型不在其中，且缓存未过期
		if len(modelCache.models) > 0 {
			modelCache.mu.Unlock()
			return fmt.Errorf("model '%s' not found in available models", model)
		}
	}
	modelCache.mu.Unlock()

	req, err := http.NewRequest("GET", "https://api.openai.com/v1/models", nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", o.apiKey))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch models: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("models API returned status %d: %s", resp.StatusCode, string(body))
	}

	// 读取响应体用于调试
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read models response body: %w", err)
	}

	var modelsResponse OpenAIModelsResponse
	if err := json.Unmarshal(body, &modelsResponse); err != nil {
		return fmt.Errorf("failed to decode models response: %w, body: %s", err, string(body))
	}

	// 更新缓存
	modelCache.mu.Lock()
	modelCache.models = make([]string, len(modelsResponse.Data))
	for i, modelData := range modelsResponse.Data {
		modelCache.models[i] = modelData.ID
	}
	modelCache.cacheTime = time.Now()
	modelCache.mu.Unlock()

	// 检查模型是否存在
	for _, modelData := range modelsResponse.Data {
		if modelData.ID == model {
			return nil
		}
	}

	return fmt.Errorf("model '%s' not found in available models", model)
}

func shouldUseNewParams(model string) bool {
	// GPT-4及更新的模型使用max_completion_tokens
	newModels := []string{"gpt-4", "gpt-4-turbo", "gpt-4o", "gpt-5", "o1-", "o3-"}
	for _, prefix := range newModels {
		if strings.HasPrefix(model, prefix) {
			return true
		}
	}
	return false
}

func (o *OpenAIProvider) GenerateCommitMessage(ctx context.Context, diff string) (string, error) {
	if o.apiKey == "" {
		return "", fmt.Errorf("openai API key is required")
	}

	// 验证模型名称
	if err := o.validateOpenAIModel(o.model); err != nil {
		return "", err
	}

	// 如果验证通过但模型看起来不常见，给出警告
	if strings.Contains(o.model, "gpt-5") {
		// 记录但不阻止，因为API验证可能通过
		// 实际使用时如果返回空内容，用户会知道模型有问题
	}

	prompt := o.template.GeneratePrompt(diff)

	request := OpenAIRequest{
		Model: o.model,
		Messages: []Message{
			{Role: "system", Content: o.template.GetSystemPrompt()},
			{Role: "user", Content: prompt},
		},
	}

	body, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", o.apiKey))

	resp, err := o.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	var responseBody []byte
	if resp.StatusCode != http.StatusOK {
		responseBody, _ = io.ReadAll(resp.Body)
		return "", fmt.Errorf("openai API returned status %d: %s (model: %s)", resp.StatusCode, string(responseBody), o.model)
	}

	// 读取响应体用于调试
	responseBody, _ = io.ReadAll(resp.Body)

	var response ChatCompletionResponse
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w, body: %s", err, string(responseBody))
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no choices in response (model: %s)", o.model)
	}

	choice := response.Choices[0]
	content := choice.Message.Content

	// 处理各种finish_reason情况
	switch choice.FinishReason {
	case "stop":
		// 正常完成
		if content == "" {
			return "", fmt.Errorf("model completed but returned empty content (model: %s)", o.model)
		}
		return content, nil
	case "length":
		// 达到token限制
		if content == "" {
			return "", fmt.Errorf("model reached token limit and returned empty content (model: %s, max_tokens: %d)", o.model, 150)
		}
		return content, nil
	case "content_filter":
		return "", fmt.Errorf("content was filtered by model (model: %s)", o.model)
	case "null":
		return "", fmt.Errorf("model response incomplete (finish_reason: null, model: %s)", o.model)
	default:
		if content == "" {
			return "", fmt.Errorf("model returned empty content with finish_reason: %s (model: %s)", choice.FinishReason, o.model)
		}
		return content, nil
	}
}

func (o *OpenAIProvider) Name() string {
	return "openai"
}
