package agents

import (
	"context"
	"fmt"
	"strings"

	"github.com/sanjbh/mailforge/internal/llm"
	"github.com/tmc/langchaingo/llms"
)

type PickerAgent struct {
	Name string
}

func NewPickerAgent() *PickerAgent {

	return &PickerAgent{
		Name: "Picker Agent",
	}
}

func (p *PickerAgent) PickBestEmail(ctx context.Context, l llms.Model, emails []string) (string, error) {

	systemPrompt, err := getSystemPrompt("picker")
	if err != nil {
		return "", fmt.Errorf("failed to get system prompt: %w", err)
	}

	var buff strings.Builder
	for i, email := range emails {
		fmt.Fprintf(&buff, "Email %d:%s\n\n", i+1, email)
	}

	res, err := llm.Generate(ctx, l, systemPrompt, buff.String(), nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate email: %w", err)
	}

	return res, nil
}
