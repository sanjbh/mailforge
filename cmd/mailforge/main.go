package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/sanjbh/mailforge/internal/agents"
	"github.com/sanjbh/mailforge/internal/config"
	"github.com/sanjbh/mailforge/internal/events"
	"github.com/sanjbh/mailforge/internal/llm"
	"github.com/sanjbh/mailforge/internal/mailer"
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

	ui.PrintInfo("Starting MailForge...")

	err = ui.StartMulti()
	if err != nil {
		log.Fatalf("Failed to start multi: %v", err)
	}

	var (
		generateEmailErrors []error
		mutex               sync.Mutex
	)

	for index, salesAgent := range salesAgents {
		wg.Add(1)

		spinner, err := ui.NewAgentSpinner(salesAgent.Name)
		if err != nil {
			log.Fatalf("Failed to create spinner: %v", err)
		}

		salesAgent.Register(spinner)

		go func(idx int, agent *agents.SalesAgent) {
			defer wg.Done()
			response, err := agent.GenerateEmail(ctx, model, prompt)
			if err != nil {
				mutex.Lock()
				generateEmailErrors = append(generateEmailErrors, err)
				mutex.Unlock()
			}
			responseEmails[idx] = response
		}(index, salesAgent)

	}
	wg.Wait()
	ui.StopMulti()

	if len(generateEmailErrors) > 0 {
		for _, e := range generateEmailErrors {
			ui.PrintFail(fmt.Sprintf("Generate email failed: %v", e))
		}
		log.Fatalf("Generate email failed with %d errors", len(generateEmailErrors))
	}

	ui.PrintInfo("Picking best email...")

	var bestEmail string
	if err = ui.RunWithSpinner("Picker Agent", func(obs events.Observer) error {
		picker := agents.NewPickerAgent()
		picker.Register(obs)

		var pickErr error
		bestEmail, pickErr = picker.PickBestEmail(ctx, model, responseEmails)
		if pickErr != nil {
			return pickErr
		}
		return nil
	}); err != nil {
		log.Fatalf("Unable to pick the best email from LLM: %v\n", err)
	}

	ui.PrintSuccess("Best email selected!")

	var mailSubject string
	if err = ui.RunWithSpinner("Subject Writer", func(obs events.Observer) error {
		writerAgent := agents.NewSubjectWriterAgent("Subject Writer")
		writerAgent.Register(obs)

		var writerErr error
		mailSubject, writerErr = writerAgent.WriteSubject(ctx, model, bestEmail)
		if writerErr != nil {
			return writerErr
		}
		return nil
	}); err != nil {
		log.Fatalf("Failed to get subject from LLM: %v\n", err)
	}

	ui.PrintSuccess("Subject generated!")

	var htmlBody string
	if err = ui.RunWithSpinner("HTML Converter", func(obs events.Observer) error {
		converterAgent := agents.NewHTMLConverterAgent("HTML Converter")
		converterAgent.Register(obs)

		var converterErr error
		htmlBody, converterErr = converterAgent.ConvertToHTML(ctx, model, bestEmail)
		if converterErr != nil {
			return converterErr
		}
		return nil

	}); err != nil {
		log.Fatalf("Unable to get body coverted to HTML by LLM: %v\n", err)
	}

	ui.PrintSuccess("Body converted to HTML!")

	mail, err := mailer.NewMail(cfg.ToEmail, cfg.FromEmail, mailSubject, htmlBody)
	if err != nil {
		log.Fatalf("Unable to create new mail. %v\n", err)
	}

	ui.PrintSuccess("Mail created!")

	var m mailer.Mailer
	if cfg.MailMock {
		m = &mailer.MockMailer{}
	} else {
		m = mailer.NewSendGridMailer(cfg)
	}

	if err := m.Send(ctx, mail); err != nil {
		log.Fatalf("Unable to send email: %v\n", err)
	}
	ui.PrintSuccess("Email sent successfully!")
}
