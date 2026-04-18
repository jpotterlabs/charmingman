# CharmingMan AI Gateway & TUI - Development Guide

## 1. Overview
The CharmingMan project consists of a unified AI Gateway (backend) and a multi-agent Chat TUI (frontend).

### Current Capabilities:
- **AI Gateway**:
    - Unified API access via `/api/v1/chat`.
    - Support for OpenAI, Anthropic, Ollama, and llama.cpp.
    - Security middleware for `X-Charming-Key` validation.
    - 30-second request timeouts for reliability.
- **Chat TUI**:
    - **Wizard State**: Guided agent configuration using the `huh` library.
    - **Dashboard State**: Multi-window management system.
    - **Window Management**: Support for dragging (mouse) and cycling focus (`Tab`).
    - **Mock Chat**: Initialized viewports for chat history display.

## 2. Setup and Launch

### Prerequisites:
- Go 1.26 or higher.
- A running instance of Ollama or llama.cpp (for local model support).

### Backend Configuration:
1. Navigate to the `backend/` directory.
2. Create or edit the `.env` file:
   ```env
   PORT=8090
   GATEWAY_API_KEY=your-secret-key-here
   OPENAI_API_KEY=sk-...
   ANTHROPIC_API_KEY=ant-...
   OLLAMA_BASE_URL=http://localhost:11434/v1
   LLAMACPP_BASE_URL=http://localhost:8081/v1
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

## 3. Manual Testing

### Testing the API (Backend):
All requests must include the `X-Charming-Key` header.
```bash
curl -X POST http://localhost:8090/api/v1/chat \
  -H "Content-Type: application/json" \
  -H "X-Charming-Key: your-secret-key-here" \
  -d '{
    "provider": "openai",
    "model": "gpt-4o",
    "prompt": "Hello!"
  }'
```

### Testing the TUI (Frontend):
1. **Wizard Flow**: Run `go run main.go`. Fill out the "Agent Name", "Model", "Persona", and "API Key". Navigate using arrow keys and Enter.
2. **Dashboard**: Once the wizard completes, you enter the dashboard.
    - **Focus**: Press `Tab` to cycle focus between windows (if multiple exist).
    - **Mouse**: You can drag windows by their borders or titles (if implemented in `manager.go`).
    - **Exit**: Press `q` or `Ctrl+C` to quit.
3. **Current Limitations**: The TUI does not yet send real requests to the backend gateway. Chat interactions are currently mock-only.

## 4. Developer Diary & Technical Notes

### Backend Decisions:
- **Unified Adapter Strategy**: Used the `fantasy` library to abstract provider differences.
- **Timeout Management**: Implemented a mandatory 30s timeout in `ProviderService.Chat`.
- **Hardened Auth**: Refactored middleware to ensure strict `X-Charming-Key` validation.
- **Error Mapping**: Mapped "provider not registered" errors to `400 Bad Request` for better API usability.

### Frontend (TUI) Decisions:
- **State Machine**: The root model manages a simple state machine (`stateWizard` -> `stateDashboard`).
- **Composition**: The `Manager` model manages multiple `Window` models, which in turn wrap content models like `ChatModel`.
- **Event Forwarding**: The TUI uses a standard Bubble Tea MVU pattern, where the root model delegates `Update` calls to the active sub-model.

---
*Note: This documentation is updated as of April 2026.*
