# CharmingMan Testing Strategy & Guide

## 1. Coverage Plan (Phases 1-3)

Phases 1-3 cover the foundation of CharmingMan: AI Gateway, SQLite persistence (sqlc), and RAG pipeline (extraction, chunking, vector search).

### AI Gateway (`backend/internal/provider`, `backend/internal/handler`)
- **Objective:** Ensure stable communication with OpenAI, Anthropic, and Local models.
- **Coverage Strategy:**
  - Mock external providers using `httptest.Server` simulating `openaicompat` endpoints to test model routing and fallback logic without real API calls.
  - Test usage logging asynchronously to SQLite via generated `sqlc` models.
  - Unit test `chat.go` to ensure prompt construction, context injection, and error handling (e.g., missing API keys) works smoothly.

### SQLite Persistence (`backend/internal/db`)
- **Objective:** Guarantee database migrations, schema generation, and query safety.
- **Coverage Strategy:**
  - **In-Memory SQLite Tests:** Use `ncruces/go-sqlite3` and `goose` in setup/teardown functions to test raw generated queries (CRUD operations on documents, chunks, usage logs).
  - Test edge cases like long document names, invalid chunk IDs, and cascade deletions.

### RAG Pipeline (`backend/internal/document`, `backend/internal/vector`)
- **Objective:** Verify document ingestion, text chunking, embedding generation, and vector retrieval.
- **Coverage Strategy:**
  - `chunker.go`: Test recursive text splitting, exact boundaries, overlapping rules, and fallback mechanisms for extremely long tokens without whitespaces.
  - `extractor.go`: Use temp directories and mock files to test text extraction (Markdown, TXT) and unsupported formats. 
  - `local.go` (Vector Store): Validate exact cosine similarity calculations, top-K selection logic, and metadata deep-copying to prevent mutation side-effects.
  - Integration Test: Write an end-to-end ingestion test from file parsing to embedding (using a Mock Embedder) to querying.

---

## 2. Testing Plan (Phases 4-5)

Phases 4-5 introduce complex TUI features and multi-agent interactions.

### Infinity Canvas & Visual Connections (`internal/tui`)
- **Objective:** Test spatial geometry, rendering limits, and component connections in Bubble Tea.
- **Testing Approach:**
  - **Headless TUI Testing:** Use `charmbracelet/bubbletea/teatest` to simulate keystrokes and assert on view models without actually rendering to a physical terminal.
  - **Coordinate Math Unit Tests:** Write tests for the Canvas model to verify viewport boundaries, panning offsets, and node coordinate translations (x, y geometry constraints).
  - **Visual Regression:** Capture string outputs of `lipgloss` rendering and use snapshot tests to ensure UI changes don't unexpectedly break the layout.

### Multi-Agent @Mentions & MCP Tool Calling
- **Objective:** Ensure accurate routing of prompts to specific agents and robust tool execution.
- **Testing Approach:**
  - **Agent Stream Testing:** Mock the `fantasy.Agent` interface to verify that `@Agent` parsing correctly segments the prompt and passes the right context.
  - **Tool Call Sandboxing:** Unit test the MCP (Model Context Protocol) executor by providing mock tools (e.g., fake CLI runner) to ensure proper JSON argument parsing, execution, and error bubbling if a tool fails.

### Whisper STT / TTS
- **Objective:** Test audio pipelines asynchronously.
- **Testing Approach:**
  - Since audio relies heavily on CGO or external binaries, abstract the STT/TTS calls behind an interface.
  - Provide dummy `.wav` byte streams and assert that the TUI transitions through `Recording -> Transcribing -> Processing` states properly.

---

## 3. Best Practices & How to Run Tests

### Running Tests
To run all backend tests:
```bash
cd backend
go test ./... -v
```

To check test coverage:
```bash
cd backend
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Mocking External Services
**1. Mocking Providers (OpenAI/Anthropic):**
Do not use real API keys in tests. Instead, spin up a local HTTP server using `httptest` and use `openaicompat` to route traffic to it. Look at `backend/internal/handler/chat_test.go` for an example of mocking an OpenAI response.

**2. Mocking Vector Stores (Pinecone):**
The project implements a local vector store (`vector.LocalStore`). In tests, inject `vector.NewLocalStore()` instead of Pinecone to keep tests fast, hermetic, and offline.

**3. Mocking Embedders:**
Implement a lightweight struct that satisfies the `vector.Embedder` interface to return deterministic float32 arrays (e.g., `[]float32{1.0, 0.0, 0.0}`) for reliable cosine similarity testing.

### Rules for Adding New Tests
1. **Idiomatic Go:** Use `github.com/stretchr/testify/assert` for clean and readable assertions.
2. **Table-Driven Tests:** For functions with multiple conditions (like `ChunkText`), use table-driven tests (`[]struct`) where feasible.
3. **No Side Effects:** File extractors and database tests must use `t.TempDir()` or in-memory databases and clean up after themselves to avoid polluting the host system.
4. **Data Race Checks:** When testing local vector stores or concurrent TUI handlers, occasionally run tests with `go test -race ./...`.
