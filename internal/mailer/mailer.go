package mailer

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"text/template"

	"github.com/go-playground/validator/v10"
)

//go:embed templates/*
var emailTemplateFS embed.FS

type Mailer interface {
	Send(ctx context.Context, mailer *Mail) error
}

type Mail struct {
	To      string `validate:"required"`
	From    string `validate:"required"`
	Subject string `validate:"required"`
	Body    string `validate:"required"`
}

func NewMail(to, from, subject, body string) (*Mail, error) {
	m := Mail{
		To:      to,
		From:    from,
		Subject: subject,
		Body:    body,
	}

	if err := validator.New().Struct(m); err != nil {
		return nil, fmt.Errorf("failed to validate mail: %w", err)
	}
	return &m, nil
}

func (m *Mail) GetCompleteEmailInHTML() (string, error) {
	tmpl, err := template.ParseFS(emailTemplateFS, "templates/email.html")
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer

	if err = tmpl.Execute(&buf, m); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
