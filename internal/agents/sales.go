package agents

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"text/template"

	"github.com/sanjbh/mailforge/internal/llm"
	"github.com/tmc/langchaingo/llms"
)

//go:embed prompts/*
var promptsFS embed.FS

type SalesAgent struct {
	Observable
	Name         string
	Instructions string
}

func NewSalesAgent(name, instructions string) *SalesAgent {
	return &SalesAgent{
		Name:         name,
		Instructions: instructions,
	}
}

func NewProfessionalSalesAgent(name string) (*SalesAgent, error) {
	systemPrompt, err := getSystemPrompt("professional")
	if err != nil {
		return nil, fmt.Errorf("failed to get system prompt: %w", err)
	}

	return &SalesAgent{
		Name:         name,
		Instructions: systemPrompt,
	}, nil
}

func NewConciseSalesAgent(name string) (*SalesAgent, error) {
	systemPrompt, err := getSystemPrompt("concise")
	if err != nil {
		return nil, fmt.Errorf("failed to get system prompt: %w", err)
	}

	return &SalesAgent{
		Name:         name,
		Instructions: systemPrompt,
	}, nil
}

func NewEngagingSalesAgent(name string) (*SalesAgent, error) {
	systemPrompt, err := getSystemPrompt("engaging")
	if err != nil {
		return nil, fmt.Errorf("failed to get system prompt: %w", err)
	}

	return &SalesAgent{
		Name:         name,
		Instructions: systemPrompt,
	}, nil
}

func getSystemPrompt(agentType string) (string, error) {

	tmpl, err := template.ParseFS(promptsFS, fmt.Sprintf("prompts/%s.txt", agentType))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var instructions bytes.Buffer

	if err := tmpl.Execute(&instructions, nil); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}
	return instructions.String(), nil

}

/* func (s *SalesAgent) GenerateEmail(
	ctx context.Context,
	l llms.Model,
	prompt string,
	streamFunc func(ctx context.Context, chunk []byte) error,
) (string, error) {
	res, err := llm.Generate(ctx, l, s.Instructions, prompt, streamFunc)
	if err != nil {
		return "", fmt.Errorf("failed to generate email: %w", err)
	}

	return res, nil

}
*/

func (s *SalesAgent) GenerateEmail(
	ctx context.Context,
	l llms.Model, prompt string,
) (string, error) {
	tokens := 0

	streamFunc := func(ctx context.Context, chunk []byte) error {
		tokens++
		s.NotifyAll(AgentEvent{Type: EventProgress, Payload: tokens})
		return nil
	}

	res, err := llm.Generate(ctx, l, s.Instructions, prompt, streamFunc)
	if err != nil {
		s.NotifyAll(AgentEvent{Type: EventError, Payload: tokens})
		return "", err
	}
	s.NotifyAll(AgentEvent{Type: EventSuccess, Payload: tokens})

	return res, nil
}
