package ui

import (
	"fmt"

	"github.com/pterm/pterm"
	"github.com/sanjbh/mailforge/internal/agents"
)

type AgentSpinner struct {
	Spinner *pterm.SpinnerPrinter
	Name    string
}

var activeMulti *pterm.MultiPrinter

func StartMulti() error {
	multiPrinter := pterm.DefaultMultiPrinter
	multi, err := multiPrinter.Start()
	if err != nil {
		return err
	}
	activeMulti = multi
	return nil
}

func StopMulti() error {
	if activeMulti != nil && activeMulti.IsActive {
		if _, err := activeMulti.Stop(); err != nil {
			return err
		}
		activeMulti = nil
	}
	return nil
}

func NewAgentSpinner(agentName string) (*AgentSpinner, error) {
	if activeMulti == nil {
		return nil, fmt.Errorf("multi printer not started, call ui.StartMulti() first")
	}
	spinner, err := pterm.DefaultSpinner.
		WithText(fmt.Sprintf("🤖 %s... 0 tokens", agentName)).
		WithWriter(activeMulti.NewWriter()).
		Start()
	if err != nil {
		return nil, fmt.Errorf("failed to start spinner: %w", err)
	}

	return &AgentSpinner{
		Spinner: spinner,
		Name:    agentName,
	}, nil
}

func (s *AgentSpinner) UpdateAgentSpinner(tokens int) {
	s.Spinner.UpdateText(fmt.Sprintf("🤖 %s... %d tokens", s.Name, tokens))
}

func (s *AgentSpinner) Success(tokens int) {
	s.Spinner.Success(fmt.Sprintf("✅ %s... %d tokens", s.Name, tokens))
}

func (s *AgentSpinner) Fail(tokens int) {
	s.Spinner.Fail(fmt.Sprintf("❌ %s... %d tokens", s.Name, tokens))
}

func PrintInfo(msg string) {
	pterm.Info.Println(msg)
}

func PrintSuccess(msg string) {
	pterm.Success.Println(msg)
}

func PrintFail(msg string) {
	pterm.Error.Println(msg)
}

func (a *AgentSpinner) Notify(event agents.AgentEvent) {
	tokens := event.Payload

	switch event.Type {
	case agents.EventProgress:
		a.UpdateAgentSpinner(tokens)
	case agents.EventSuccess:
		a.Success(tokens)
	case agents.EventError:
		a.Fail(tokens)
	}
}
