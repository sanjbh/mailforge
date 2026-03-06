package agents

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"
)

//go:embed prompts/*.txt
var promptsFS embed.FS

func getSystemPrompt(agentType string) (string, error) {
	tmpl, err := template.ParseFS(promptsFS, fmt.Sprintf("prompts/%s.txt", agentType))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var instructions bytes.Buffer

	if err := tmpl.Execute(&instructions, nil); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}
	return instructions.String(), nil
}
