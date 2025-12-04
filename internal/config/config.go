package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	APIKeys  map[string]string `mapstructure:"api_keys"`
	Model    string            `mapstructure:"model"`
	Provider string            `mapstructure:"provider"`
	Editor   string            `mapstructure:"editor"`
	Custom   CustomConfig      `mapstructure:"custom"`
}

type CustomConfig struct {
	URL    string `mapstructure:"url"`
	APIKey string `mapstructure:"api_key"`
	Model  string `mapstructure:"model"`
}

func Load() (*Config, error) {
	viper.SetConfigName("aicommit")
	viper.SetConfigType("yaml")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	// Only read from home directory .config
	configDir := filepath.Join(homeDir, ".config", "aicommit")
	viper.AddConfigPath(configDir)

	// Set defaults
	viper.SetDefault("model", "claude-3-sonnet-20240229")
	viper.SetDefault("provider", "claude")
	viper.SetDefault("api_keys", map[string]string{
		"claude":   "",
		"openai":   "",
		"deepseek": "",
	})

	viper.SetEnvPrefix("AICOMMIT")
	viper.AutomaticEnv()

	// Try to read config file, but don't error if it doesn't exist
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

func (c *Config) GetAPIKey(provider string) string {
	// Environment variables take precedence
	envKey := fmt.Sprintf("AICOMMIT_%s_API_KEY", provider)
	// Check env var using both standard GetEnv and Viper (which handles case sensitivity and prefixes)
	if val := os.Getenv(envKey); val != "" {
		return val
	}

	if apiKey, ok := c.APIKeys[provider]; ok && apiKey != "" {
		return apiKey
	}
	return ""
}
