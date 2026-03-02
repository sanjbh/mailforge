package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	LLMKey      string `envconfig:"LLM_KEY" default:"123457678897987979"`
	LLMBaseURL  string `envconfig:"LLM_BASE_URL" default:"http://localhost:11434/v1"`
	LLMModel    string `envconfig:"LLM_MODEL" default:"qwen3:8b-q4_K_M"`
	MailMock    bool   `envconfig:"MAIL_MOCK" default:"true"`
	FromEmail   string `envconfig:"FROM_EMAIL" default:"Test <test@example.com>"`
	ToEmail     string `envconfig:"TO_EMAIL" default:"Recipient <recipient@example.com>"`
	SendGridKey string `envconfig:"SENDGRID_KEY" default:"SG.123457678897987979"`
}

func Load() (*Config, error) {
	var cfg Config

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return &cfg, nil
}
