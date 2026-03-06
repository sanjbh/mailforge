# mailforge

A multi-agent AI pipeline written in Go that autonomously drafts, evaluates, formats, and delivers cold sales emails. Built as a production-quality Go implementation of an agentic workflow, using [langchaingo](https://github.com/tmc/langchaingo) and [Ollama](https://ollama.com/) for local LLM inference.

---

## What It Does

mailforge orchestrates a pipeline of AI agents that work together to produce and send a polished cold sales email:

1. **Three Sales Agents** run concurrently — Professional, Engaging, and Concise — each drafting a unique cold email
2. **Picker Agent** evaluates all three drafts and selects the best one
3. **Subject Writer Agent** generates a compelling subject line for the winning email
4. **HTML Converter Agent** transforms the plain text/markdown email into beautifully styled HTML
5. **Mailer** delivers the final email — either via SendGrid or a local mock that writes `output.html`

---

## Pipeline Flow

```
┌─────────────────────────────────────────────────────────┐
│                   SALES AGENTS (concurrent)             │
│  ┌──────────────┐ ┌──────────────┐ ┌─────────────────┐  │
│  │ Professional │ │   Engaging   │ │     Concise     │  │
│  │ Sales Agent  │ │ Sales Agent  │ │  Sales Agent    │  │
│  └──────┬───────┘ └──────┬───────┘ └────────┬────────┘  │
└─────────┼────────────────┼──────────────────┼───────────┘
          │                │                  │
          └────────────────┼──────────────────┘
                           │  3 email drafts
                           ▼
                   ┌───────────────┐
                   │  Picker Agent │  selects best email
                   └───────┬───────┘
                           │  winning draft
              ┌────────────┴────────────┐
              ▼                         ▼
   ┌──────────────────┐     ┌───────────────────────┐
   │  Subject Writer  │     │    HTML Converter      │
   │     Agent        │     │       Agent            │
   └────────┬─────────┘     └───────────┬────────────┘
            │ subject                   │ html body
            └────────────┬──────────────┘
                         ▼
                  ┌─────────────┐
                  │    Mailer   │  MockMailer or SendGrid
                  └─────────────┘
```

---

## Architecture

```
cmd/
└── mailforge/
    └── main.go              # Entry point and pipeline orchestration

internal/
├── config/
│   └── config.go            # Environment-based configuration via envconfig
├── events/
│   └── events.go            # Observer interface, AgentEvent, EventType
├── llm/
│   └── client.go            # LLM client with streaming + variadic options
├── agents/
│   ├── observable.go        # Observable struct — embed in any agent
│   ├── prompts.go           # Shared prompt loader using embed.FS
│   ├── sales.go             # Professional, Engaging, Concise sales agents
│   ├── picker.go            # Picker agent — selects best email
│   └── formatter.go         # Subject writer + HTML converter agents
├── mailer/
│   ├── mailer.go            # Mailer interface + Mail struct + validation
│   ├── mock.go              # MockMailer — writes output.html
│   ├── sendgrid.go          # SendGridMailer — real email delivery
│   └── templates/
│       └── email.html       # HTML email template
└── ui/
    └── ui.go                # Terminal UI — spinners via pterm
```

---

## Design Patterns

### Observer Pattern
Agents are **Observables** — they emit events as they work. The terminal spinner is an **Observer** — it watches agents and updates the UI in real time.

```
SalesAgent (Observable)  →  emits EventProgress, EventSuccess, EventError
      ↓  Register()
AgentSpinner (Observer)  →  reacts by updating the terminal spinner
```

This fully decouples agents from the UI. Agents never import or know about `pterm` or any UI library. You could swap the spinner for a file logger, a metrics collector, or a WebSocket reporter without touching a single line of agent code.

### Dependency Injection
The `Mailer` interface allows plugging in different delivery backends:

```go
type Mailer interface {
    Send(ctx context.Context, mail *Mail) error
}
```

Controlled via the `MAIL_MOCK` environment variable:
- `MAIL_MOCK=true` → writes `output.html` locally
- `MAIL_MOCK=false` → sends via SendGrid

### Concurrent Agent Execution
The three sales agents run in parallel using goroutines and `sync.WaitGroup`, with errors collected safely via `sync.Mutex`:

```go
for index, salesAgent := range salesAgents {
    wg.Add(1)
    go func(idx int, agent *agents.SalesAgent) {
        defer wg.Done()
        response, err := agent.GenerateEmail(ctx, model, prompt)
        ...
    }(index, salesAgent)
}
wg.Wait()
```

### RunWithSpinner
Single agents (Picker, Subject Writer, HTML Converter) use a `RunWithSpinner` helper that encapsulates the spinner lifecycle — `StartMulti`, spinner creation, and `StopMulti` via `defer` — keeping `main.go` clean:

```go
var bestEmail string
if err = ui.RunWithSpinner("Picker Agent", func(obs events.Observer) error {
    picker := agents.NewPickerAgent()
    picker.Register(obs)
    var pickErr error
    bestEmail, pickErr = picker.PickBestEmail(ctx, model, responseEmails)
    return pickErr
}); err != nil {
    log.Fatalf("Unable to pick best email: %v", err)
}
```

### Variadic LLM Options
`llm.Generate` accepts variadic `llms.CallOption` parameters, making it easy to pass provider-specific options without breaking existing callers:

```go
func Generate(ctx, llm, system, prompt, streamFunc, opts ...llms.CallOption) (string, error)

// Example — cap tokens for subject writer:
llm.Generate(ctx, l, systemPrompt, emailBody, streamFunc, llms.WithMaxTokens(50))
```

### Separation of Concerns via `internal/events`
The `Observer` interface and `AgentEvent` type live in a neutral `internal/events` package. This prevents circular imports between `agents` and `ui` — both import `events`, but never each other:

```
agents  →  events  ←  ui
```

---

## Tech Stack

| Component | Library |
|---|---|
| LLM Abstraction | [langchaingo](https://github.com/tmc/langchaingo) |
| Local LLM | [Ollama](https://ollama.com/) + `qwen3:8b-q4_K_M` |
| Email Delivery | [SendGrid Go SDK](https://github.com/sendgrid/sendgrid-go) |
| Config | [envconfig](https://github.com/kelseyhightower/envconfig) |
| Validation | [go-playground/validator](https://github.com/go-playground/validator) |
| Terminal UI | [pterm](https://github.com/pterm/pterm) |
| Templating | Go standard `text/template` + `embed.FS` |

---

## Prerequisites

- Go 1.22+
- [Ollama](https://ollama.com/) running locally
- `qwen3:8b-q4_K_M` model pulled

```bash
ollama pull qwen3:8b-q4_K_M
ollama serve
```

---

## Installation

```bash
git clone https://github.com/sanjbh/mailforge.git
cd mailforge
go mod download
```

---

## Configuration

Create a `.env` file in the project root:

```env
# LLM
LLM_KEY=dummy
LLM_BASE_URL=http://localhost:11434/v1
LLM_MODEL=qwen3:8b-q4_K_M

# Mailer
MAIL_MOCK=true
FROM_EMAIL=Your Name <you@example.com>
TO_EMAIL=Recipient Name <recipient@example.com>

# SendGrid (only needed if MAIL_MOCK=false)
SENDGRID_KEY=SG.your-key-here
```

---

## Usage

```bash
go run ./cmd/mailforge
```

If `MAIL_MOCK=true`, the final email is written to `output.html` in the project root. Open it in a browser to see the rendered result.

To send a real email via SendGrid, set `MAIL_MOCK=false` and provide a valid `SENDGRID_KEY`.

---

## Terminal Output

```
ℹ Starting MailForge...
✅ Concise Sales Agent... 528 tokens
✅ Engaging Sales Agent... 786 tokens
✅ Professional Sales Agent... 940 tokens
ℹ Picking best email...
✅ Picker Agent... 901 tokens
✅ Best email selected!
✅ Subject Writer... 42 tokens
✅ Subject generated!
✅ HTML Converter... 1162 tokens
✅ Body converted to HTML!
✅ Mail created!
✅ Email sent successfully!
```

---

## Extending mailforge

Adding a new agent is three steps:

1. Create a struct embedding `Observable` and write a method that calls `llm.Generate` with a streaming callback that fires `NotifyAll`
2. Add a prompt file under `internal/agents/prompts/`
3. Wire it up in `main.go` with `RunWithSpinner` — the UI wiring never touches agent code:

```go
ui.RunWithSpinner("My New Agent", func(obs events.Observer) error {
    agent := agents.NewMyAgent("My New Agent")
    agent.Register(obs)
    result, err = agent.DoWork(ctx, model, input)
    return err
})
```

---

## License

MIT
