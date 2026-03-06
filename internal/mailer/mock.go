package mailer

import (
	"context"
	"fmt"
	"os"
)

type MockMailer struct{}

func (m *MockMailer) Send(ctx context.Context, mailer_m *Mail) error {

	file, err := os.Create("output.html")
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	body, err := mailer_m.GetCompleteEmailInHTML()
	if err != nil {
		return fmt.Errorf("failed to get email body: %w", err)
	}

	_, err = file.WriteString(body)
	if err != nil {
		return fmt.Errorf("failed to write to output file: %w", err)
	}

	return nil
}
