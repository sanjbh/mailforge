package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/sanjbh/mailforge/internal/agents"
	"github.com/sanjbh/mailforge/internal/config"
	"github.com/sanjbh/mailforge/internal/llm"
	"github.com/sanjbh/mailforge/internal/ui"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Config failed: %v", err)
	}

	ctx := context.Background()
	var wg sync.WaitGroup

	salesAgents := make([]*agents.SalesAgent, 0)

	conciseSalesAgent, err := agents.NewConciseSalesAgent("Concise Sales Agent")
	if err != nil {
		log.Fatalf("Failed to create concise sales agent: %v", err)
	}

	engagingSalesAgent, err := agents.NewEngagingSalesAgent("Engaging Sales Agent")
	if err != nil {
		log.Fatalf("Failed to create engaging sales agent: %v", err)
	}

	professionalSalesAgent, err := agents.NewProfessionalSalesAgent("Professional Sales Agent")
	if err != nil {
		log.Fatalf("Failed to create professional sales agent: %v", err)
	}

	salesAgents = append(salesAgents, conciseSalesAgent, engagingSalesAgent, professionalSalesAgent)

	prompt := "Write a cold sales email addressed to 'Dear CEO'"

	model, err := llm.New(cfg)
	if err != nil {
		log.Fatalf("LLM failed: %v", err)
	}

	responseEmails := make([]string, len(salesAgents))

	fmt.Println()

	multi, err := ui.StartMulti()
	if err != nil {
		log.Fatalf("Failed to start multi: %v", err)
	}

	for index, salesAgent := range salesAgents {
		wg.Add(1)

		spinner, err := ui.NewAgentSpinner(multi, salesAgent.Name)
		if err != nil {
			log.Fatalf("Failed to create agent spinner: %v", err)
		}

		go func(agent *agents.SalesAgent, idx int, s *ui.AgentSpinner) {
			defer wg.Done()

			var tokens int

			streamCallback := func(ctx context.Context, chunk []byte) error {
				tokens++
				s.UpdateAgentSpinner(tokens)
				return nil
			}

			response, err := agent.GenerateEmail(ctx, model, prompt, streamCallback)
			if err != nil {
				s.Fail(tokens)
			} else {
				s.Success(tokens)
			}
			responseEmails[idx] = response
		}(salesAgent, index, spinner)
	}
	wg.Wait()
	multi.Stop()

	picker := agents.NewPickerAgent()

	ui.PrintInfo("Picking best email...")
	bestEmail, err := picker.PickBestEmail(ctx, model, responseEmails)
	if err != nil {
		log.Fatalf("Pick best email failed: %v", err)
	}

	ui.PrintSuccess(fmt.Sprintf("Best email: %s\n", bestEmail))
}
