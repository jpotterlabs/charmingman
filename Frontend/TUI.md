# CharmingMan – Multi‑Agent ChatTUI  
*A terminal‑native, modular, Go‑powered AI assistant built with the Charmbracelet ecosystem.*

---

## 1. Executive Overview

| Feature | What it does | Why it matters |
|---------|--------------|----------------|
| **Multi‑pane chat** | Run several agents side‑by‑side (GPT‑4, Claude, local Ollama, etc.) | Compare answers, run parallel reasoning, reduce hallucinations |
| **Persona engine** | Assign a “personality” (style, tone, core values) to each agent | Keeps agents consistent & memorable |
| **Document “Stage”** | Render PDFs, code, markdown, and diagrams in‑terminal | One‑stop workspace for all artifacts |
| **Tool‑calling & RAG** | Built‑in Whisper, TTS, database lookups, external APIs | Makes the agent truly “agentic” |
| **Stateful chatrooms** | Store conversation history on disk or DB | Avoids sending long context on every request |
| **Extensible middleware** | Swap between OpenAI, Anthropic, local, Hugging‑Face in one call | Future‑proof & cost‑efficient |
| **Mouse‑friendly split panes** | Resize, collapse, and focus panes with mouse or keyboard | IDE‑style UX in the terminal |

---

## 2. Layered Architecture

```
┌─────────────────────-----------──┐
│   1. Frontend TUI (This Document)│
└───────────────────────-----------┘
        ▲
        │
┌──────────────────────---┐
│   2. Web API AI Gateway │
└─────────────────────----┘
        ▲
        │
┌──────────────────────-------─┐
│   3. Backend Document Engine│
└──────────────────────-------─┘
```

### 2.1 Frontend (TUI Layer)

| Sub‑layer | Libraries | Responsibility |
|-----------|-----------|----------------|
| **Orchestrator** | `bubbletea` | MVU loop, view switching (Home / Wizard / Chat / Document) |
| **Wizard** | `huh` | Guided agent creation (name, model, persona, tools, API keys) |
| **UI Components** | `bubbles` | Text input, progress bars, viewports, menus |
| **Styling** | `lipgloss` | Flexbox‑style layout, borders, adaptive spacing |
| **Rendering** | `glow` | Markdown → ANSI (tables, code blocks, diagrams) |
| **Theme Sync** | `catwalk` | Detect terminal theme, apply accent colors for providers |

> **State machine** (in the TUI)  
> `Home → Wizard → Chat → Document → ...`  
> Each view emits `tea.Cmd`s that are routed to the middleware.

### 2.2 AI Gateway (Middleware)

| Component | Purpose | Example API |
|-----------|---------|-------------|
| **Fantasy (Agentic framework)** | Define agents, protocols, tool calls, structured responses | `fantasy.NewAgent("summarizer")` |
| **Provider Translator** | Route internal request to OpenAI, Anthropic, Ollama, local HF models | `Translate(ctx, req) → *ProviderRequest` |
| **Streaming Engine** | SSE/WS for token‑by‑token updates | `OpenAIStreaming(req) → chan Token` |
| **Structured Response Parser** | Map JSON‑schema to Go structs | `json.Unmarshal(msg, &Response{})` |
| **Model Cache** | Download & store local models (GGUF, Safetensors) | `huggingface.Download(modelID)` |

> **Middleware Flow**  
> `TUI → Msg → Agentic Engine → Provider → LLM → Response → TUI`

### 2.3 Backend (Document & Multi‑Modal Engine)

| Responsibility | Tools / Libraries |
|----------------|-------------------|
| **RAG** | Vector store (ChromaDB/FAISS), embedding model (OpenAI/Local) |
| **PDF / Document** | `unidoc/unipdf` → text → chunk → embed |
| **Audio** | Whisper (STT) via `portaudio`/`ffmpeg`, TTS via ElevenLabs/OpenAI |
| **Telephony** | Twilio / SignalWire webhook handling |
| **Tool Calling** | Local Go functions wrapped as JSON‑schema endpoints |
| **MCP / WebSocket** | Agent‑to‑agent communication, file & DB access |
| **Persistence** | SQLite/Badger for chatrooms, state, and lore |
| **Background Workers** | Goroutines for long‑running jobs, SSE to TUI |

---

## 3. UI Layout & Navigation

### 3.1 Dashboard (Multi‑pane “Beam”)

```
+---------------------------------------------+
|   Sidebar (Agent Library, MCP, Telemetry)   |
+---------------------------------------------+
|  [Chat #1]  [Chat #2]  [Chat #3]  [Chat #4] |
+---------------------------------------------+
|            Stage (Glow) – Documents & Code   |
+---------------------------------------------+
|  Input Buffer (shared for all chats)         |
+---------------------------------------------+
```

- **Resizable Panes** – Mouse or `Ctrl+Arrows`  
- **Active Focus** – Bold/thicker border + accent color  
- **Keybindings**  

| Key | Action | Component |
|-----|--------|-----------|
| `Tab` | Cycle focus | `bubbletea` |
| `Ctrl+N` | Open Wizard | `huh` |
| `Ctrl+S` | Toggle Stage | `lipgloss` |
| `/` | Command mode (tool call, model switch) | `bubbles` |
| `Ctrl+B` | Collapse sidebar | `catwalk` |

### 3.2 Sidebar – “Intelligence Hub”

| Section | UI Element | Description |
|---------|------------|-------------|
| Agent Library | Vertical list | Each entry shows icon (`☁️`, `🏠`) and name |
| MCP Servers | Status badges | E.g., “Filesystem: Active” |
| Token Telemetry | Sparkline | Real‑time token cost per session |

### 3.3 Stage – Document Engine

- **Document Preview** – `glow` renders PDF, markdown, Mermaid, PlantUML.  
- **Code & Artifacts** – Full‑width view, syntax highlighted.  
- **Thinking Log** – Collapsible drawer showing chain‑of‑thought, tool calls.

---

## 4. Agent & Persona Design

### 4.1 Personality Struct

```go
type Personality struct {
    Name      string   `json:"name"`
    Archetype string   `json:"archetype"` // e.g., "Stoic Mentor"
    Traits    []string `json:"traits"`    // e.g., ["verbose", "witty"]
    Speech    string   `json:"speech"`    // e.g., "Pirate slang"
    Values    []string `json:"values"`    // e.g., ["honesty", "patience"]
}
```

### 4.2 Wizard Integration

| Wizard Step | UI Control | What it Sets |
|-------------|------------|--------------|
| Name | `huh.Input` | Agent identifier |
| Model | `huh.Select` | Provider + model ID |
| Personality | `huh.MultiSelect`, `huh.Select` | Archetype, traits, speech |
| Tools | `huh.MultiSelect` | List of allowed JSON‑schema tools |
| API Keys | `huh.Password` | Secure entry, stored encrypted |

### 4.3 System Prompt Template

```go
const systemPromptFmt = `# Roleplay Protocol: %s
You are %s, a %s.
## Constraints:
- Tone: %s
- Linguistic Style: %s
- Core Values: %s

## Interaction Rules:
Never break character. If asked about your AI nature, respond in the voice of %s.`
```

Inject this into the LLM request as the system role. The `Fantasy` framework will then enforce the structured response schema.

---

## 5. State Machine (High‑Level Flow)

```
Idle
 └─► Wizard (Agent Creation)
      └─► Chat (Bubble Tea)
          ├─► AwaitingInput
          ├─► Processing (spinner)
          ├─► StreamingResponse
          └─► Document (Glow)
```

Each state is a separate `bubbletea` model; transitions are triggered by `tea.Cmd`s such as `StartChat`, `OpenDocument`, `FinishWizard`, etc.

---

## 6. Implementation Roadmap

| Phase | Deliverables | Key Libraries |
|-------|--------------|---------------|
| **I – Foundations** | Project scaffolding, TUI skeleton | `bubbletea`, `lipgloss` |
| **II – Wizard & Agent Model** | `huh` wizard, Persona struct | `huh`, `fantasy` |
| **III – Middleware** | Provider translator, streaming, structured responses | `fantasy`, OpenAI SDK, Anthropic SDK |
| **IV – Document Engine** | PDF → text, embedding, vector store | `unidoc`, `chroma`, `go-embed` |
| **V – Audio & Telephony** | Whisper, TTS, Twilio integration | `portaudio`, ElevenLabs API |
| **VI – UI Polish** | Resizable panes, adaptive theme, status bar | `catwalk`, `bubbles` |
| **VII – Extras** | Lorebook, role‑play logs, persistence, tests | SQLite, `go-memdb`, `testify` |

---

## 7. Example Code Snippets

### 7.1 Bubble Tea Model – Chat Pane

```go
type ChatPane struct {
    ID          string
    AgentID     string
    History     []Message
    Input       string
    StreamCh    chan string
    err         error
}

func (m *ChatPane) Init() tea.Cmd {
    return nil
}

func (m *ChatPane) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.Type == tea.KeyEnter {
            return m, startChatRequest(m)
        }
    case tokenMsg:
        m.History = append(m.History, msg)
        return m, nil
    case errorMsg:
        m.err = msg
        return m, nil
    }
    return m, nil
}
```

### 7.2 Fantasy Agent Definition

```go
func NewSummarizerAgent() *fantasy.Agent {
    return fantasy.NewAgent("summarizer",
        fantasy.WithTools(searchTool, translateTool),
        fantasy.WithPersona(Personality{
            Name:      "Summarizer",
            Archetype: "Concise Analyst",
            Traits:    []string{"brief", "direct"},
            Speech:    "neutral",
            Values:    []string{"accuracy"},
        }),
    )
}
```

### 7.3 Provider Translator

```go
func Translate(req *AgentRequest) (*ProviderRequest, error) {
    switch req.Provider {
    case "openai":
        return &ProviderRequest{
            Endpoint: openaiEndpoint,
            Model:    req.Model,
            Tokens:   req.Tokens,
        }, nil
    case "anthropic":
        return &ProviderRequest{
            Endpoint: anthropicEndpoint,
            Model:    req.Model,
            Tokens:   req.Tokens,
        }, nil
    // …
    }
}
```

---

## 8. Testing & Quality

| Test Area | Strategy |
|-----------|----------|
| **Unit** | `go test` + `testify` for models, provider translator, RAG pipeline |
| **Integration** | `go test -run TestIntegration` – start a local OpenAI mock server, send a full chat request, verify streaming |
| **End‑to‑End** | `go test -run TestTUI` – use `bubbletea`’s `tea.NewProgram` with a fake UI to simulate key presses |
| **Performance** | Benchmark vector store queries, token streaming latency |

---

## 9. Packaging & Distribution

- **CLI** – `charmingman` binary (single static binary, Go 1.22+)
- **Registry** – `charmingman install <pkg>` fetches a pre‑built UI theme or agent package from a GitHub repo
- **Docker** – optional container for running as a service (e.g., in a CI pipeline)
- **Configuration** – `.env.local` for API keys, `config.yaml` for default provider and theme

---

## 10. Final Thoughts

CharmingMan is not just a chat client; it’s a **terminal‑native AI workspace**:

- **Modular** – Add or replace a provider, a tool, or a storage backend without touching the UI.  
- **Extensible** – Plug in new agents, add new personas, or connect to custom APIs.  
- **Performance‑friendly** – Streaming, background workers, and local models keep the UI snappy.  
- **Professional look** – Catwalk + Lipgloss + Glow give it the polish of a web dashboard while staying fully terminal‑based.

With this blueprint, the next step is to bootstrap the skeleton, wire up the wizard, and iterate on the first chat experience. Happy building!
