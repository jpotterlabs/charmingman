# CharmingMan 🕶️

**CharmingMan** is a sophisticated, multi-agent ChatTUI built exclusively on the **Charm** ecosystem. It provides a high-performance terminal interface for interacting with various LLM providers (OpenAI, Anthropic, Ollama, llama.cpp) and features a powerful RAG (Retrieval-Augmented Generation) pipeline and a revolutionary **Infinity Canvas**.

## ✨ Features

- **🧠 Multi-Agent Orchestration**: Create and manage multiple AI agents with distinct personas. Agents share a persistent **RoomID** context and can be routed using **@mentions**.
- **🌌 Infinity Canvas**: A spatial workspace for your agents. Navigate a non-linear graph of thoughts with camera panning, zooming, and world-to-screen coordinate mapping.
- **📚 Knowledge & RAG**: Ground your agents in your own data. Upload and index PDFs, Markdown, and Text files with automated context injection.
- **🖼️ "The Stage"**: A dedicated document preview window for inspecting RAG sources and AI-generated artifacts using `Glow`.
- **🪟 YAML-Driven Layouts**: Define your workspace in `layout.yaml`. Features semantic validation, auto-rescaling for different terminal sizes, and draggable/resizable windows.
- **🔌 Unified AI Gateway**: A backend service that abstracts provider differences, offering a single API with built-in security like prompt redaction and path-traversal protection.
- **🛠️ Wizard-driven Setup**: Easily configure new agents and providers through an interactive `Huh?` wizard.

## 🚀 Quick Start

### Prerequisites
- Go 1.26 or higher.
- (Optional) [Ollama](https://ollama.ai/) for local model support.
- (Optional) [Pinecone](https://www.pinecone.io/) account for managed vector storage.

### Installation
1. Clone the repository:
   ```bash
   git clone git@github.com:OWNER/REPO.git
   cd charmingman
   ```
2. Install dependencies:
   ```bash
   go mod tidy
   cd backend && go mod tidy
   ```

### Configuration
Create a `.env` file in the `backend/` directory:
```env
PORT=8090
OPENAI_API_KEY=sk-...
ANTHROPIC_API_KEY=ant-...
PINECONE_API_KEY=your-pinecone-key
PINECONE_INDEX=your-index-name
DOCUMENTS_ROOT=./documents
```

### Running
1. **Start the AI Gateway**:
   ```bash
   cd backend
   go run cmd/gateway/main.go
   ```
2. **Launch the TUI** (in a new terminal):
   ```bash
   go run main.go
   ```

## 🛠️ Architecture

CharmingMan is built with a layered architecture:
- **Frontend**: `Bubble Tea`, `Lipgloss`, `Huh`, `Glow`, and a custom **Spatial Canvas Engine**.
- **Middleware (AI Gateway)**: Go-based service using the `fantasy` library for provider abstraction and multi-agent routing.
- **Backend (Intelligence Engine)**: SQLite for persistence, Pinecone/Local for vector storage, and a custom RAG pipeline with deep-copy mutation safety.

## 🗺️ Roadmap
- [x] Phase 1: AI Gateway & TUI Architecture
- [x] Phase 2: Persistence & Local SQLite
- [x] Phase 3: Intelligence & Knowledge (RAG)
- [x] Phase 4: Advanced TUI & Infinity Canvas (Swarms, @Mentions, Layouts)
- [ ] Phase 5: Voice & Multimedia (Whisper/TTS/MCP)

---
Built with ❤️ using [Charm](https://charm.sh/).
