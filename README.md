# MailForge 🤖✉️

MailForge is a professional **Agentic Workflow** CLI tool designed to generate and pick the best cold sales emails using a multi-agent system. It features a decoupled, event-driven UI built on the **Observer Pattern**.

## 🚀 Features

- **Agentic Multi-Agent Workflow**: Orchestrates a collaboration between specialized generation agents and an evaluator agent.
- **Specialized LLM Agents**: Features Concise, Engaging, and Professional agents that run in parallel.
- **AI-Powered Evaluator**: A dedicated Picker agent acts as the "manager" to analyze all generated outputs and select the most effective email.
- **Dynamic UI**: Real-time progress tracking with a multi-spinner terminal interface using [pterm](https://github.com/pterm/pterm).
- **Extensible Architecture**: Easily add new agents with minimal UI plumbing thanks to a generic Observer-based notification system.

## 🏗️ Architecture

The project follows a clean, modular structure:

- `cmd/mailforge/`: Application entry point and concurrency orchestration.
- `internal/agents/`: Core agent logic and the Observer pattern infrastructure.
- `internal/ui/`: State-managed UI components and spinner encapsulation.
- `internal/llm/`: Simplified wrapper for LLM interactions via LangChainGo.
- `internal/config/`: Environment-based configuration management.

## 🛰️ The Observer Pattern

We use the Observer pattern to decouple agents from the terminal UI. Agents emit generic events, and the UI reacts to them without the agents knowing the UI exists.

### How it works:
1. **Subject (Agents)**: Any agent can embed `Observable` to broadcast events (`EventProgress`, `EventSuccess`, `EventError`).
2. **Observer (Spinner)**: The `AgentSpinner` implements the `Observer` interface.
3. **Decoupling**: To link them, you simply call `agent.Register(spinner)`. The `main.go` file stays clean, focused only on execution logic.

## 🛠️ Usage

### Prerequisites
- Go 1.21+
- An OpenAI-compatible API key (configured in `.env`)

### Setup
1. Clone the repository.
2. Copy `.env.example` to `.env` and add your LLM credentials.
3. Install dependencies:
   ```bash
   go mod tidy
   ```

### Running the tool
```bash
go run cmd/mailforge/main.go
```

## 📝 Example Output
```text
🤖 Concise Sales Agent... 42 tokens
🤖 Engaging Sales Agent... 58 tokens
🤖 Professional Sales Agent... 35 tokens
✅ Picking best email...
```

## 📄 License
This project is licensed under the MIT License.
