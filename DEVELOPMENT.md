# CharmingMan AI Gateway & TUI - Development Guide

## 1. Overview
The CharmingMan project consists of a unified AI Gateway (backend) and a multi-agent Chat TUI (frontend).

### Current Capabilities:
- **AI Gateway**:
    - Unified API access via `/api/v1/chat`.
    - Support for OpenAI, Anthropic, Ollama, and llama.cpp.
    - **Security**: Prompt redaction in logs and path-traversal protections for document uploads.
    - **RAG Implementation**: Intelligent retrieval from uploaded documents with context injection.
- **Chat TUI**:
    - **Infinity Canvas**: Spatial workspace with camera panning/zooming and world-to-screen coordinate mapping.
    - **Multi-Agent Orchestration**: @mention routing and shared RoomID context.
    - **YAML Layouts**: Configurable `layout.yaml` with semantic validation and auto-rescaling.

## 2. Phase 4: Advanced TUI & Infinity Canvas

Phase 4 transformed CharmingMan into a spatial workspace, enabling complex agentic orchestration.

### 🌌 Infinity Canvas
The Infinity Canvas is a non-linear spatial environment where windows are placed in a global "world" coordinate system.
- **Spatial Layout**: Windows are no longer restricted to a rigid grid but exist at $(x, y)$ coordinates in a theoretically infinite plane.
- **Camera System**: Includes panning (moving the viewpoint) and zooming (scaling the view).
- **Coordinate Mapping**: The engine maps `World Coordinates` (where a window is in the infinite space) to `Screen Coordinates` (where it appears in your terminal window).

### 🧠 Multi-Agent Orchestration
Agents now operate within a collaborative environment.
- **@Mention Routing**: Users can direct prompts to specific agents using `@AgentName`. The TUI parses these mentions and routes the message to the corresponding agent model.
- **Persistent RoomID**: Every session generates a `RoomID`. All agents in that session share this ID, allowing the AI Gateway to maintain a shared conversational context across multiple participants.

### 🪟 YAML Layouts (`layout.yaml`)
Workspaces are now defined using a declarative YAML schema.
- **Schema**: Defines panes, their initial positions, sizes, and linked components (e.g., `ChatFeed`, `DocViewer`).
- **Semantic Validation**: The loader verifies the layout for logical errors (e.g., overlapping fixed windows or invalid component IDs).
- **Auto-Rescaling**: When the terminal window is resized, the layout engine automatically adjusts window proportions and coordinates to maintain visual consistency.

## 3. Phase 3: Intelligence & Knowledge (RAG)

Phase 3 focused on grounding agents with your own documents.

### 🧩 RAG Pipeline Components
1. **Document Extractor**: Supports `.txt`, `.md`, and `.pdf` files.
2. **Chunker**: Splits large documents into smaller pieces with context overlap.
3. **Embedder**: Generates vector embeddings for each chunk.
4. **Vector Store**: Supports Local (in-memory) and Pinecone (cloud) storage.

## 4. Security & Safety

- **Prompt Redaction**: Sensitive information in logs is automatically redacted before being persisted to the database.
- **Path-Traversal Protection**: The document ingestion service strictly validates file paths to prevent unauthorized file access via `../` patterns.
- **Deep-Copy Mutation Safety**: The internal data structures used during RAG retrieval are deep-copied to prevent unexpected state changes across concurrent requests.

## 5. Setup and Launch

### Prerequisites:
- Go 1.26 or higher.
- A running instance of Ollama or llama.cpp (for local model support).

### Configuration:
1. Navigate to the `backend/` directory.
2. Create or edit the `.env` file:
   ```env
   PORT=8090
   GATEWAY_API_KEY=your-secret-key-here
   OPENAI_API_KEY=sk-... # REQUIRED for embeddings
   ANTHROPIC_API_KEY=ant-...
   
   # Pinecone Configuration (Optional, defaults to LocalStore if missing)
   PINECONE_API_KEY=your-pinecone-key
   PINECONE_INDEX=your-index-name
   
   # Document Storage
   DOCUMENTS_ROOT=./documents
   ```

### Running the Project:
- **Start the Gateway**:
  ```bash
  cd backend
  go run cmd/gateway/main.go
  ```
- **Start the TUI**:
  ```bash
  go run main.go
  ```

## 6. Next Steps (Roadmap)
- **Phase 5: Voice & Multimedia (Whisper/TTS/MCP)**:
    - Implement the Model Context Protocol (MCP) for local tool calling.
    - Enable Whisper STT for voice input.
    - Add Text-to-Speech (TTS) for agent responses.

---
*Note: This documentation is updated as of June 2026.*
