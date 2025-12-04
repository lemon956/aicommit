package model

import (
	"fmt"

	"github.com/aicommit/aicommit/internal/config"
)

func NewProvider(cfg *config.Config) (Provider, error) {
	switch cfg.Provider {
	case "claude":
		return NewClaudeProvider(cfg.GetAPIKey("claude"), cfg.Model), nil
	case "openai":
		return NewOpenAIProvider(cfg.GetAPIKey("openai"), cfg.Model), nil
	case "deepseek":
		return NewDeepSeekProvider(cfg.GetAPIKey("deepseek"), cfg.Model), nil
	case "custom":
		return NewCustomProvider(cfg.Custom.URL, cfg.Custom.APIKey, cfg.Custom.Model), nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", cfg.Provider)
	}
}
