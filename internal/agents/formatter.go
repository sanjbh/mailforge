package agents

import (
	"context"
	"fmt"

	"github.com/sanjbh/mailforge/internal/events"
	"github.com/sanjbh/mailforge/internal/llm"
	"github.com/tmc/langchaingo/llms"
)

type SubjectWriterAgent struct {
	Name string
	Observable
}

func NewSubjectWriterAgent(name string) *SubjectWriterAgent {
	return &SubjectWriterAgent{
		Name: name,
	}
}

func (s *SubjectWriterAgent) WriteSubject(ctx context.Context, l llms.Model, emailBody string) (string, error) {
	systemPrompt, err := getSystemPrompt("subject_writer")
	if err != nil {
		return "", fmt.Errorf("failed to get system prompt: %w", err)
	}

	tokens := 0
	streamFunc := func(ctx context.Context, chunks []byte) error {
		tokens++
		s.NotifyAll(events.AgentEvent{
			Type:    events.EventProgress,
			Payload: tokens,
		})
		return nil
	}

	res, err := llm.Generate(ctx, l, systemPrompt, emailBody, streamFunc, llms.WithMaxTokens(50))
	if err != nil {
		s.NotifyAll(events.AgentEvent{
			Type:    events.EventError,
			Payload: tokens,
		})
		return "", fmt.Errorf("failed to generate subject: %w", err)
	}
	s.NotifyAll(events.AgentEvent{
		Type:    events.EventSuccess,
		Payload: tokens,
	})

	return res, nil
}

type HTMLConverterAgent struct {
	Name string
	Observable
}

func NewHTMLConverterAgent(name string) *HTMLConverterAgent {
	return &HTMLConverterAgent{
		Name: name,
	}
}

func (h *HTMLConverterAgent) ConvertToHTML(ctx context.Context, l llms.Model, emailBody string) (string, error) {
	systemPrompt, err := getSystemPrompt("html_converter")
	if err != nil {
		return "", fmt.Errorf("failed to get system prompt: %w", err)
	}

	tokens := 0

	streamFunc := func(ctx context.Context, chunk []byte) error {
		tokens++
		h.NotifyAll(events.AgentEvent{Type: events.EventProgress, Payload: tokens})
		return nil
	}

	response, err := llm.Generate(ctx, l, systemPrompt, emailBody, streamFunc)
	if err != nil {
		h.NotifyAll(events.AgentEvent{Type: events.EventError, Payload: tokens})
		return "", fmt.Errorf("failed to convert to HTML: %w", err)
	}
	h.NotifyAll(events.AgentEvent{Type: events.EventSuccess, Payload: tokens})

	return response, nil

}
