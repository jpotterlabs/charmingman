# CharmingMan AI Gateway & TUI - Development Guide

## 1. Overview
The CharmingMan project consists of a unified AI Gateway (backend) and a multi-agent Chat TUI (frontend).

### Current Capabilities:
- **AI Gateway**:
    - Unified API access via `/api/v1/chat`.
    - **Multimodal Transcription**: `/api/v1/transcribe` endpoint for Whisper STT.
    - Support for OpenAI, Anthropic, Ollama, and llama.cpp.
    - **Security**: Prompt redaction in logs and path-traversal protections for document uploads.
    - **RAG Implementation**: Intelligent retrieval from uploaded documents with context injection.
- **Chat TUI**:
    - **Infinity Canvas**: Spatial workspace with camera panning/zooming and world-to-screen coordinate mapping.
    - **Multi-Agent Orchestration**: @mention routing and shared RoomID context.
    - **Voice Input**: Integrated `sox` recording triggered by the 'v' key for hands-free interaction.
    - **YAML Layouts**: Configurable `layout.yaml` with semantic validation and auto-rescaling.

## 2. Phase 5: Voice & Multimodal Interaction

Phase 5 introduced speech-to-text capabilities, allowing users to interact with agents using their voice.

### 🎙️ Whisper STT Integration
The AI Gateway now features a transcription handler that interfaces with OpenAI's Whisper model.
- **Endpoint**: `POST /api/v1/transcribe`
- **Logic**: Receives a `.wav` file (multipart form), sends it to Whisper, and returns the transcribed text.
- **Requirement**: `OPENAI_API_KEY` must be configured in the gateway.

### ⌨️ TUI Voice Trigger
The TUI includes a `VoiceInputModel` that captures audio directly from the terminal.
- **Trigger**: Pressing the 'v' key starts a recording session.
- **Recording Engine**: Uses the system's `rec` command (from `sox`) to capture audio for 3 seconds.
- **Routing**: Once transcribed, the text is automatically sent to the primary chat agent via the same routing logic used for text input.

## 3. Stability & Resilience Fixes (Phase 5)

Several key stability improvements were implemented during Phase 5:
- **Bounded Chat History**: The Gateway now limits history retrieval to the last 10 messages to prevent token overflow and ensure low-latency responses.
- **Log Redaction**: All request/response logging in the database is sanitized. Personal identifiable information or sensitive prompt content is redacted to maintain user privacy.
- **Coordinate Scaling**: Fixed a bug in the Infinity Canvas where window coordinates would drift on terminal resize. The new scaling logic ensures that world-to-screen mapping remains pixel-perfect regardless of terminal dimensions.

## 4. Phase 4: Advanced TUI & Infinity Canvas

Phase 4 transformed CharmingMan into a spatial workspace.

### 🌌 Infinity Canvas
The Infinity Canvas is a non-linear spatial environment where windows are placed in a global "world" coordinate system.
- **Spatial Layout**: Windows are no longer restricted to a rigid grid.
- **Camera System**: Includes panning and zooming.
- **Coordinate Mapping**: Maps `World Coordinates` to `Screen Coordinates`.

### 🧠 Multi-Agent Orchestration
- **@Mention Routing**: Route prompts to specific agents using `@AgentName`.
- **Persistent RoomID**: Shared conversational context across multiple participants.

## 5. Phase 3: Intelligence & Knowledge (RAG)

Phase 3 focused on grounding agents with your own documents.

### 🧩 RAG Pipeline Components
1. **Document Extractor**: Supports `.txt`, `.md`, and `.pdf` files.
2. **Chunker**: Splits large documents into smaller pieces.
3. **Embedder**: Generates vector embeddings.
4. **Vector Store**: Supports Local and Pinecone storage.

## 6. Security & Safety

- **Prompt Redaction**: Sensitive information in logs is automatically redacted.
- **Path-Traversal Protection**: Strict validation of file paths for document ingestion.
- **Deep-Copy Mutation Safety**: Data structures used during RAG retrieval are deep-copied to prevent state crosstalk.

## 7. Setup and Launch

### Prerequisites:
- Go 1.26 or higher.
- `sox` (specifically `rec`) for voice input.
- A running instance of Ollama or llama.cpp (for local model support).

### Configuration:
1. Navigate to the `backend/` directory.
2. Create or edit the `.env` file:
   ```env
   PORT=8090
   GATEWAY_API_KEY=your-secret-key-here
   OPENAI_API_KEY=sk-... # REQUIRED for embeddings and Whisper STT
   ANTHROPIC_API_KEY=ant-...
   
   # Pinecone Configuration (Optional)
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

## 8. Next Steps (Roadmap)
- **Phase 6: TTS & MCP (Local Tool Calling)**:
    - Implement the Model Context Protocol (MCP).
    - Add Text-to-Speech (TTS) for agent responses.

---
*Note: This documentation is updated as of June 2026.*
