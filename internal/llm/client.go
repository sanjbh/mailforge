package llm

import (
	"context"
	"fmt"

	"github.com/sanjbh/mailforge/internal/config"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

func New(cfg *config.Config) (llms.Model, error) {
	return openai.New(
		openai.WithBaseURL(cfg.LLMBaseURL),
		openai.WithModel(cfg.LLMModel),
		openai.WithToken(cfg.LLMKey),
	)
}

func Generate(ctx context.Context, llm llms.Model, system, prompt string) (string, error) {
	msgs := []llms.MessageContent{
		{
			Role: llms.ChatMessageTypeSystem,
			Parts: []llms.ContentPart{
				llms.TextContent{
					Text: system,
				},
			},
		},
		{
			Role: llms.ChatMessageTypeHuman,
			Parts: []llms.ContentPart{
				llms.TextContent{
					Text: prompt,
				},
			},
		},
	}

	// log.Printf("Generating content with system: %s, prompt: %s", system, prompt)

	res, err := llm.GenerateContent(ctx, msgs)
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	if len(res.Choices) == 0 {
		return "", fmt.Errorf("no choices returned")
	}

	if res.Choices[0].Content == "" {
		return "", fmt.Errorf("empty response content")
	}

	return res.Choices[0].Content, nil
}
