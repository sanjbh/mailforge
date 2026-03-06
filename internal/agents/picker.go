package agents

import (
	"context"
	"fmt"
	"strings"

	"github.com/sanjbh/mailforge/internal/events"
	"github.com/sanjbh/mailforge/internal/llm"
	"github.com/tmc/langchaingo/llms"
)

type PickerAgent struct {
	Name string
	Observable
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

	tokens := 0

	streamFunc := func(context context.Context, chunk []byte) error {
		tokens++
		p.NotifyAll(events.AgentEvent{Type: events.EventProgress, Payload: tokens})
		return nil
	}

	res, err := llm.Generate(ctx, l, systemPrompt, buff.String(), streamFunc)
	if err != nil {
		p.NotifyAll(events.AgentEvent{Type: events.EventError, Payload: tokens})
		return "", err
		// return "", fmt.Errorf("failed to generate email: %w", err)
	}

	p.NotifyAll(events.AgentEvent{Type: events.EventSuccess, Payload: tokens})

	return res, nil
}
