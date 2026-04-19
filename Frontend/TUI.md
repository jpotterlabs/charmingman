# CharmingMan – Multi‑Agent ChatTUI  
*A terminal‑native, modular, Go‑powered AI assistant built with the Charmbracelet ecosystem.*

---

## 1. Executive Overview

| Feature | What it does | Why it matters |
|---------|--------------|----------------|
| **Multi‑pane chat** | Run several agents side‑by‑side (GPT‑4, Claude, local Ollama, etc.) | Compare answers, run parallel reasoning, reduce hallucinations |
| **Voice Input (v5)** | Hands-free prompt entry via Whisper STT | Interaction becomes seamless and "Star Trek" like |
| **Infinity Canvas** | A spatial workspace with infinite panning and zooming | Non-linear thought mapping and agent organization |
| **Persona engine** | Assign a “personality” (style, tone, core values) to each agent | Keeps agents consistent & memorable |
| **Document “Stage”** | Render PDFs, code, markdown, and diagrams in‑terminal | One‑stop workspace for all artifacts |
| **Tool‑calling & RAG** | Built‑in Whisper, TTS, database lookups, external APIs | Makes the agent truly “agentic” |
| **Stateful chatrooms** | Store conversation history on disk or DB | Bounded to last 10 messages for speed |
| **Extensible middleware** | Swap between providers (OpenAI, Anthropic, Ollama) | Future‑proof & cost‑efficient |

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
| **Voice Suite** | `sox`, `whisper` | Capture audio via 'v' key and send to gateway |
| **Wizard** | `huh` | Guided agent creation (name, model, persona, tools, API keys) |
| **UI Components** | `bubbles` | Text input, progress bars, viewports, menus |
| **Styling** | `lipgloss` | Flexbox‑style layout, borders, adaptive spacing |
| **Rendering** | `glow` | Markdown → ANSI (tables, code blocks, diagrams) |
| **Theme Sync** | `catwalk` | Detect terminal theme, apply accent colors |

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

### 3.2 Keybindings

| Key | Action | Component |
|-----|--------|-----------|
| `Tab` | Cycle focus | `bubbletea` |
| **`v`** | **Start Voice Input (Phase 5)** | `VoiceInputModel` |
| `Ctrl+N` | Open Wizard | `huh` |
| `Ctrl+S` | Toggle Stage | `lipgloss` |
| `/` | Command mode (tool call, model switch) | `bubbles` |
| `Ctrl+B` | Collapse sidebar | `catwalk` |
| `+/-` | Zoom in/out (Canvas) | `SpatialEngine` |
| `WASD` | Pan camera (Canvas) | `SpatialEngine` |

---

## 4. Implementation Roadmap

| Phase | Deliverables | Status |
|-------|--------------|--------|
| **I – Foundations** | Project scaffolding, TUI skeleton | COMPLETED |
| **II – Wizard & Agent Model** | `huh` wizard, Persona struct | COMPLETED |
| **III – Middleware** | Provider translator, streaming, RAG | COMPLETED |
| **IV – Spatial Canvas** | Infinity Canvas, Panning, Zooming | COMPLETED |
| **V – Voice & Multimodal** | Whisper STT, `sox` integration | COMPLETED |
| **VI – TTS & MCP** | Text-to-Speech, Local tool calling | IN PROGRESS |

---

## 5. Final Thoughts

CharmingMan is not just a chat client; it’s a **terminal‑native AI workspace**:

- **Multimodal** – Interact via text or voice.
- **Spatial** – Organize agents in an infinite workspace.
- **Performance‑friendly** – Bounded history and local models keep the UI snappy.
- **Professional look** – Catwalk + Lipgloss + Glow polish.

Happy building!
