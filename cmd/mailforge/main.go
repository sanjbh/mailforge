package main

import (
	"context"
	"fmt"
	"log"
	"os"
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

	err = ui.StartMulti()
	if err != nil {
		log.Fatalf("Failed to start multi: %v", err)
	}

	pickerSpinner, err := ui.NewAgentSpinner("Picker Agent")
	picker := agents.NewPickerAgent()
	picker.Register(pickerSpinner)

	for index, salesAgent := range salesAgents {
		wg.Add(1)

		spinner, _ := ui.NewAgentSpinner(salesAgent.Name)
		salesAgent.Register(spinner)

		go func(idx int, agent agents.SalesAgent) {
			defer wg.Done()
			response, err := agent.GenerateEmail(ctx, model, prompt)
			if err != nil {
				ui.PrintFail(err.Error())
				os.Exit(1)
			}
			responseEmails[idx] = response
		}(index, *salesAgent)

	}
	wg.Wait()

	ui.PrintInfo("Picking best email...")
	bestEmail, err := picker.PickBestEmail(ctx, model, responseEmails)
	if err != nil {
		log.Fatalf("Pick best email failed: %v", err)
	}
	ui.StopMulti()

	ui.PrintSuccess(fmt.Sprintf("Best email: %s\n", bestEmail))
}
