# CharmingMan - Development Guide & Diary

## 1. Overview
CharmingMan is a multi-agent ChatTUI built on the Charm ecosystem. It consists of a unified AI Gateway (backend) and a spatial Infinity Canvas (frontend).

## 2. Technical Architecture

### Spatial Manager & Coordinate Mapping
The Infinity Canvas uses a "World Coordinate" system. The mapping from global world positions to terminal grid cells is handled in `internal/tui/manager.go`:
```go
screenX := int(float64(window.X - OffsetX) * Zoom)
screenY := int(float64(window.Y - OffsetY) * Zoom)
```
This allows for smooth camera panning and zooming across an arbitrary plane.

### Multi-Agent Swarms & Shared Context
Agents are orchestrated by the `Manager` using a broadcast message pattern (`RouteMsg`). 
- **Persistence**: All agents in a workspace share a `RoomID`. 
- **Context Injection**: The AI Gateway automatically retrieves room history and injects it into LLM prompts, ensuring that different agents in the same room are aware of each other's responses.

### MCP Tooling Client
CharmingMan implements a custom Model Context Protocol (MCP) client in `backend/internal/mcp/client.go`. 
- **Transport**: JSON-RPC over `stdin/stdout`.
- **Discovery**: The Gateway spawns MCP server processes defined in `MCP_SERVERS` and maps their tools into the `fantasy.AgentTool` interface.

## 3. Implementation History (Phases 1-5)

### Phase 1: AI Gateway Foundation
- Unified provider interface using the `fantasy` library.
- Initial Gin router and authentication middleware.

### Phase 2: Persistence & Workflow
- SQLite integration with `sqlc`.
- Automated task management with `Taskfile.yaml`.
- Secure Agent defined DTOs.

### Phase 3: Knowledge Base (RAG)
- Document ingestion pipeline (PDF/Text extraction).
- Chunking strategy with context overlap.
- Vector store integration (Pinecone/Local).

### Phase 4: Infinity Canvas
- Spatial layout engine.
- @mention routing system.
- Shared room history persistence.

### Phase 5: Multimedia & Tools
- Whisper STT integration (Voice Input).
- OpenAI TTS integration (Agent Speech).
- Model Context Protocol (MCP) support for local tools.

## 4. Setup & Launch

### Prerequisites
- Go 1.26+
- `sox` (system audio utility)
- OpenAI API Key

### Configuration
1. Navigate to `backend/` and create a `.env` file.
2. Configure your keys and server command:
```env
PORT=8090
GATEWAY_API_KEY=your-secret-key
OPENAI_API_KEY=sk-...
MCP_SERVERS="npx @modelcontextprotocol/server-filesystem /path/to/docs"
```

### Building & Testing
```bash
task build-backend
task build-tui
task test
```

## 5. Next Steps (Roadmap)
- **Phase 6**: Cloud Sync & Encrypted Backups.
- **Phase 7**: E2E Terminal Testing with `vhs`.
- **Phase 8**: Refined Authorization & Human-in-the-Loop for sensitive tools.
