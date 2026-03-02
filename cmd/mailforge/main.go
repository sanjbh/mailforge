package main

import (
	"context"
	"log"
	"sync"

	"github.com/sanjbh/mailforge/internal/agents"
	"github.com/sanjbh/mailforge/internal/config"
	"github.com/sanjbh/mailforge/internal/llm"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Config failed: %v", err)
	}

	ctx := context.Background()

	salesAgents := make([]*agents.SalesAgent, 0)

	conciseSalesAgent, err := agents.NewConciseSalesAgent("Concise Sales Agent")
	if err != nil {
		log.Fatalf("Failed to create concise sales agent: %v", err)
	}
	salesAgents = append(salesAgents, conciseSalesAgent)

	engagingSalesAgent, err := agents.NewEngagingSalesAgent("Engaging Sales Agent")
	if err != nil {
		log.Fatalf("Failed to create engaging sales agent: %v", err)
	}
	salesAgents = append(salesAgents, engagingSalesAgent)

	professionalSalesAgent, err := agents.NewProfessionalSalesAgent("Professional Sales Agent")
	if err != nil {
		log.Fatalf("Failed to create professional sales agent: %v", err)
	}
	salesAgents = append(salesAgents, professionalSalesAgent)

	var wg sync.WaitGroup

	prompt := "Write a cold sales email addressed to 'Dear CEO'"

	model, err := llm.New(cfg)
	if err != nil {
		log.Fatalf("LLM failed: %v", err)
	}

	responseEmails := make([]string, len(salesAgents))

	log.Println("Generating emails...")

	for index, salesAgent := range salesAgents {
		wg.Add(1)
		go func(agent *agents.SalesAgent, idx int) {
			defer wg.Done()
			log.Printf("Generating email for %s\n", agent.Name)
			response, err := agent.GenerateEmail(ctx, model, prompt)
			if err != nil {
				log.Printf("Generate email failed for %s: %v\n", agent.Name, err)
			}
			// log.Printf("Response for %s: %s\n", agent.Name, response)
			responseEmails[idx] = response
		}(salesAgent, index)
	}
	wg.Wait()
	log.Println("All emails generated")

	picker := agents.NewPickerAgent()

	log.Println("Picking best email...")
	bestEmail, err := picker.PickBestEmail(ctx, model, responseEmails)
	if err != nil {
		log.Fatalf("Pick best email failed: %v", err)
	}

	log.Printf("Best email: %s\n", bestEmail)
}
