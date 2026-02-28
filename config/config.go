package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	TelegramToken string `json:"telegram_token"`
	OllamaModel   string `json:"ollama_model"`
	SystemPrompt  string `json:"system_prompt"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		OllamaModel:  "llama3.2",
		SystemPrompt: "You are a helpful Telegram chatbot. Keep responses short and conversational.",
	}

	data, err := os.ReadFile("config.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read config.json: %w", err)
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config.json: %w", err)
	}

	if cfg.TelegramToken == "" {
		return nil, fmt.Errorf("telegram_token is required")
	}

	return cfg, nil
}
