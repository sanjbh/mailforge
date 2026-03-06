package mailer

import (
	"context"
	"fmt"

	"github.com/sanjbh/mailforge/internal/config"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailer struct {
	config *config.Config
}

func NewSendGridMailer(cfg *config.Config) *SendGridMailer {
	return &SendGridMailer{
		config: cfg,
	}
}

func (s *SendGridMailer) Send(ctx context.Context, m *Mail) error {

	client := sendgrid.NewSendClient(s.config.SendGridKey)

	from := mail.NewEmail("", m.From)
	to := mail.NewEmail("", m.To)

	body, err := m.GetCompleteEmailInHTML()
	if err != nil {
		return fmt.Errorf("failed to get email body: %w", err)
	}

	message := mail.NewSingleEmail(from, m.Subject, to, "", body)

	res, err := client.Send(message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	if res.StatusCode >= 400 {
		return fmt.Errorf("sendgrid returned status %d: %s", res.StatusCode, res.Body)
	}
	return nil
}
