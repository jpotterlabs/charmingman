# CharmingMan Features & Capabilities

**CharmingMan** is built to be a robust, multi-agent terminal assistant with deep integrations and a cutting-edge interface.

## 🚀 Key Features

### 1. Multi-Agent Swarms
- **Collaborative Orchestration**: Agents operate within a shared **RoomID** context, allowing them to remain aware of the conversation history.
- **Direct Routing**: Use **@mentions** (e.g., `@Architect`, `@Coder`) to route prompts to specific agent models.
- **Thinking Drawer**: Visualize the internal chain-of-thought and reasoning process of each agent as it "thinks."

### 2. Infinity Canvas
- **Spatial TUI**: A non-linear workspace where windows can be placed anywhere on an infinite plane.
- **Dynamic Interaction**: Full mouse support for dragging and resizing windows, with camera panning and zooming for global navigation.
- **World-to-Screen Mapping**: Intelligent rendering that maps high-resolution world coordinates to terminal character grids.

### 3. Voice & Multimodal (New in Phase 5)
- **Whisper STT**: Integrated Speech-to-Text via OpenAI Whisper for hands-free interaction.
- **TUI Recording**: Trigger voice recording directly from the terminal using the 'v' key.
- **Audio Capture**: Leverages system-level `sox` for high-fidelity audio capture in the TUI.

### 4. Intelligence & RAG
- **Deep Document Integration**: Ground agent responses in PDFs, Markdown, and Text files.
- **Automated Context Injection**: Relevant chunks of your data are intelligently injected into prompts based on the query.
- **RAG Safety**: Deep-copy mutation protections ensure that concurrent RAG queries do not interfere with agent states.

### 5. Advanced TUI Engine
- **YAML Layouts**: Fully configurable workspaces defined via `layout.yaml`. Supports semantic validation and auto-rescaling for various terminal sizes.
- **The Stage**: A dedicated high-performance viewport using `Glow` for inspecting documents, code, and artifacts.
- **Adaptive Themes**: Visual consistency across all terminal themes using `Catwalk`.

### 6. AI Gateway (Middleware)
- **Unified Provider API**: Single endpoint access for OpenAI, Anthropic, Ollama, and llama.cpp.
- **Multimodal API**: Dedicated transcription endpoint for audio-to-text processing.
- **Security-First Design**: Includes prompt redaction in persistent logs and strict path-traversal protections.

## 🗺️ Implementation Roadmap

### Phase 1: AI Gateway & TUI Architecture [COMPLETED]
- Project scaffolding and basic Bubble Tea loop.
- Initial provider translation layer.

### Phase 2: Persistence & Local SQLite [COMPLETED]
- Migration-based database management with `goose` and `sqlc`.
- Chat history and agent state persistence.

### Phase 3: Intelligence & Knowledge (RAG) [COMPLETED]
- Document extraction and chunking pipeline.
- Local and cloud vector storage.

### Phase 4: Advanced TUI & Infinity Canvas [COMPLETED]
- Spatial canvas engine with panning and zooming.
- Multi-agent @mention routing and RoomID context.
- YAML-driven layout system with semantic validation.

### Phase 5: Voice & Multimedia [COMPLETED]
- Whisper STT (Speech-to-Text) integration.
- Bounded chat history and stability improvements.
- Secure log redaction and coordinate scaling fixes.

### Phase 6: TTS & MCP [IN PROGRESS]
- Model Context Protocol (MCP) for local tool calling.
- Text-to-Speech (TTS) for agent responses.
- Telephony features and multi-modal agentic triggers.
