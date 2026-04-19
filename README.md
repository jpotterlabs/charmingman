# CharmingMan 🕶️

**CharmingMan** is a sophisticated, multi-agent ChatTUI built exclusively on the **Charm** ecosystem. It provides a high-performance terminal interface for interacting with various LLM providers (OpenAI, Anthropic, Ollama, llama.cpp) and features a powerful spatial "Infinity Canvas" for orchestrating agent swarms.

## ✨ Features

- **🧠 [Multi-Agent Orchestration](./docs/features/MULTI_AGENT.md)**: Create and manage multiple AI agents with distinct personas. Use `@mentions` to route prompts and coordinate reasoning across your swarm.
- **🌌 [Infinity Canvas](./docs/features/INFINITY_CANVAS.md)**: A dynamic, spatial layout engine. Pan and zoom across an infinite coordinate plane to organize your workspace.
- **📚 [Knowledge & RAG](./docs/features/KNOWLEDGE_RAG.md)**: Ground your agents in your own data. Upload and index PDFs, Markdown, and Text files with automated chunking and vector storage (Pinecone/Local).
- **🎙️ [Multimodal Suite](./docs/features/MULTIMODAL.md)**: Trigger voice recording in the TUI to prompt agents and listen to their responses via high-quality speech synthesis.
- **🛠️ [MCP Tooling](./docs/features/MCP_TOOLS.md)**: Extend your agents with local capabilities like filesystem access, shell execution, and more via the Model Context Protocol.
- **🔌 [Unified AI Gateway](./docs/features/AI_GATEWAY.md)**: A backend service that abstracts provider differences, offering a single API for all your models with persistent room history.

## 🚀 Quick Start

### Prerequisites
- Go 1.26 or higher.
- `sox` (required for voice features): `brew install sox` (macOS) or `sudo apt install sox` (Linux).
- (Optional) [Ollama](https://ollama.ai/) for local model support.
- (Optional) [Pinecone](https://www.pinecone.io/) account for managed vector storage.

### Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/jpotterlabs/charmingman.git
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
GATEWAY_API_KEY=charming-secret-key
OPENAI_API_KEY=sk-...
ANTHROPIC_API_KEY=ant-...
PINECONE_API_KEY=your-pinecone-key
PINECONE_INDEX=your-index-name
DOCUMENTS_ROOT=./documents
MCP_SERVERS="npx @modelcontextprotocol/server-filesystem /path/to/docs"
```

### Running
1. **Start the AI Gateway**:
   ```bash
   task run-backend
   ```
2. **Launch the TUI** (in a new terminal):
   ```bash
   task run-tui
   ```

## 🛠️ Architecture

CharmingMan is built with a layered architecture:
- **Frontend**: `Bubble Tea v2`, `Lipgloss v2`, `Huh v2`. Features a custom compositor for spatial rendering.
- **Middleware (AI Gateway)**: Go-based service using the `fantasy` library for provider abstraction.
- **Backend (Intelligence Engine)**: SQLite for persistence, Pinecone/Local for vector storage, and a custom document processing pipeline.

## 📖 Documentation

For detailed information on each system, see our feature guides:
- [AI Gateway Guide](./docs/features/AI_GATEWAY.md)
- [Multi-Agent Guide](./docs/features/MULTI_AGENT.md)
- [Infinity Canvas Guide](./docs/features/INFINITY_CANVAS.md)
- [Knowledge & RAG Guide](./docs/features/KNOWLEDGE_RAG.md)
- [Multimodal Guide](./docs/features/MULTIMODAL.md)
- [MCP Tools Guide](./docs/features/MCP_TOOLS.md)
- [Testing Strategy](./TESTING.md)
- [Developer Diary](./DEVELOPMENT.md)

---
Built with ❤️ using [Charm](https://charm.sh/).
