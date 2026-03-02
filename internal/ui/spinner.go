package ui

import (
	"fmt"

	"github.com/pterm/pterm"
)

type AgentSpinner struct {
	Spinner *pterm.SpinnerPrinter
	Name    string
}

func StartMulti() (*pterm.MultiPrinter, error) {
	multi := pterm.DefaultMultiPrinter
	return multi.Start()
}

func NewAgentSpinner(p *pterm.MultiPrinter, agentName string) (*AgentSpinner, error) {
	spinner, err := pterm.DefaultSpinner.
		WithText(fmt.Sprintf("🤖 %s... 0 tokens", agentName)).
		WithWriter(p.NewWriter()).
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
	s.Spinner.Success(fmt.Sprintf("🤖 %s... %d tokens", s.Name, tokens))
}

func (s *AgentSpinner) Fail(tokens int) {
	s.Spinner.Fail(fmt.Sprintf("🤖 %s... %d tokens", s.Name, tokens))
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
