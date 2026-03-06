package ui

import (
	"fmt"

	"github.com/pterm/pterm"
	"github.com/sanjbh/mailforge/internal/events"
)

type AgentSpinner struct {
	Spinner *pterm.SpinnerPrinter
	Name    string
}

var activeMulti *pterm.MultiPrinter

func StartMulti() error {
	pterm.DefaultSpinner = *pterm.DefaultSpinner.WithSequence("⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏").
		WithStyle(pterm.NewStyle(pterm.FgCyan))
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

func (a *AgentSpinner) Notify(event events.AgentEvent) {
	tokens := event.Payload

	switch event.Type {
	case events.EventProgress:
		a.UpdateAgentSpinner(tokens)
	case events.EventSuccess:
		a.Success(tokens)
	case events.EventError:
		a.Fail(tokens)
	}
}

func RunWithSpinner(name string, fn func(obs events.Observer) error) error {
	defer StopMulti()
	if err := StartMulti(); err != nil {
		return err
	}

	spinner, err := NewAgentSpinner(name)
	if err != nil {
		return err
	}

	return fn(spinner)
}
