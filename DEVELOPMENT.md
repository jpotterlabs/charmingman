# CharmingMan AI Gateway & TUI - Development Guide

## 1. Overview
The CharmingMan project consists of a unified AI Gateway (backend) and a multi-agent Chat TUI (frontend).

### Current Capabilities:
- **AI Gateway**:
    - Unified API access via `/api/v1/chat`.
    - Support for OpenAI, Anthropic, Ollama, and llama.cpp.
    - Security middleware for `X-Charming-Key` validation.
    - **RAG Implementation**: Intelligent retrieval from uploaded documents.
- **Chat TUI**:
    - **Wizard State**: Guided agent configuration using the `huh` library.
    - **Dashboard State**: Multi-window management system.
    - **Window Management**: Draggable windows (mouse support) and focus cycling (`Tab`).
    - **"The Stage"**: A dedicated window for document preview and RAG source inspection.

## 2. Phase 3: Intelligence & Knowledge (RAG)

Phase 3 focused on grounding agents with your own documents. This involved a significant backend expansion and new TUI components.

### 🧩 RAG Pipeline Components
1. **Document Extractor**: Supports `.txt`, `.md`, and `.pdf` files. PDF extraction is powered by `github.com/ledongthuc/pdf`.
2. **Chunker**: Splits large documents into smaller pieces (1000 characters) with a 200-character overlap for context continuity.
3. **Embedder**: Generates vector embeddings for each chunk (currently using `OpenAI text-embedding-3-small`).
4. **Vector Store**:
    - **LocalStore**: In-memory storage for development.
    - **PineconeStore**: Production-grade managed vector database.

### 📦 New Dependencies
- `github.com/pinecone-io/go-pinecone`: Go SDK for Pinecone.
- `github.com/ledongthuc/pdf`: Native Go PDF parser for text extraction.
- `github.com/gin-gonic/gin`: HTTP web framework for the AI Gateway.
- `github.com/google/uuid`: Used for document and chunk identification.

## 3. Setup and Launch

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

## 4. Manual Testing

### Testing the API (Backend):
All requests must include the `X-Charming-Key` header.
```bash
# Upload a document
curl -X POST http://localhost:8090/api/v1/documents \
  -H "Content-Type: application/json" \
  -H "X-Charming-Key: your-secret-key-here" \
  -d '{
    "title": "My Knowledge",
    "path": "path/to/my/document.pdf"
  }'

# Chat with RAG
curl -X POST http://localhost:8090/api/v1/chat \
  -H "Content-Type: application/json" \
  -H "X-Charming-Key: your-secret-key-here" \
  -d '{
    "provider": "openai",
    "model": "gpt-4o",
    "prompt": "What does my knowledge say about X?",
    "use_rag": true
  }'
```

### Testing the TUI (Frontend):
1. **Wizard Flow**: Run `go run main.go`. Fill out the "Agent Name", "Model", etc.
2. **Dashboard**:
    - **RAG Support**: Toggle `Use RAG` in the wizard.
    - **"The Stage"**: Inspect the "The Stage" window to see the initial knowledge base message.
    - **Window Management**: Click and drag window title bars or borders to move/resize windows. Press `Tab` to cycle focus.

## 5. Developer Diary & Technical Notes

### Backend Decisions:
- **Transactional Ingestion**: Implemented compensation logic in `DocumentService.AddDocument` to ensure clean rollbacks (deleting vectors/DB records) if ingestion fails halfway.
- **Provider Abstraction**: Extended the provider system to support a unified `VectorStore` interface, allowing easy switching between local in-memory and cloud-based stores.

### Frontend (TUI) Decisions:
- **Window Manager**: Introduced a new `Manager` model in `internal/tui` to handle the complexity of multiple, overlapping windows.
- **Mouse Integration**: Built custom drag/resize logic into the `Manager` to make the TUI feel like a modern workspace while staying terminal-native.
- **Document Model**: Created a `DocumentModel` wrapped in a `Window` (conceptualized as "The Stage") to display long-form content separate from chat history.

## 6. Next Steps (Roadmap)
- **Phase 4: Multi-Agent Swarms & Tools**:
    - Implement the Model Context Protocol (MCP) for local tool calling.
    - Enable agents to communicate with each other (Agent-to-Agent).
    - Add a "Thinking" drawer to visualize agent reasoning chains.

---
*Note: This documentation is updated as of April 2026.*
