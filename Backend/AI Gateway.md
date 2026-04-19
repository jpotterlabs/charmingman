Build a unified web api that functions as a AI Gateway, providing unified API access to many AI models through a single endpoint.

# AI Gateway

The AI Gateway provides a unified API to access hundreds of AI models through a single endpoint, with built-in budgets, usage monitoring, fallbacks, and security.

---

## 🚀 Key Endpoints

### 1. Chat Completions
`POST /api/v1/chat`

Handles multi-agent routing, history retrieval, and RAG injection.
- **Provider Support**: OpenAI, Anthropic, Ollama, llama.cpp.
- **Features**: @mention routing, RoomID persistence, Bounded history (last 10 messages).

### 2. Transcription (STT) - *New in Phase 5*
`POST /api/v1/transcribe`

Provides high-fidelity speech-to-text using OpenAI's Whisper model.
- **Request Type**: `multipart/form-data`
- **Fields**:
  - `file`: The audio file to transcribe (e.g., `.wav`, `.mp3`).
- **Response**:
  ```json
  {
    "text": "The transcribed text from the audio."
  }
  ```
- **Requirements**: Requires a valid `OPENAI_API_KEY` on the gateway server.

### 3. Document Ingestion
`POST /api/v1/documents`

Upload and index documents for RAG.
- **Supported Formats**: `.pdf`, `.md`, `.txt`.
- **Security**: Built-in path-traversal protection.

---

## 🔒 Security & Performance

- **Prompt Redaction**: Sensitive information is automatically redacted from database logs.
- **History Bounding**: History retrieval is capped at the last 10 messages to ensure performance and stay within model context limits.
- **Deep-Copy Mutation Safety**: State objects are deep-copied during routing to prevent crosstalk between concurrent agent requests.
- **Unified Auth**: Access all downstream providers using a single `GATEWAY_API_KEY` (X-Charming-Key).

---

## 🛠️ Configuration (Example .env)

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

---

## 🗺️ Capabilities

**Text Generation**
- Reasoning & Problem Solving
- Multi-agent Swarm Routing
- Persistent Context (RoomID)

**Audio (Phase 5)**
- **STT**: Whisper-1 integration for voice input.
- **TTS**: (Planned) ElevenLabs / OpenAI TTS integration.

**Document Management**
- Automated Chunking & Embedding
- Local & Cloud (Pinecone) Vector Stores
- PDF Extraction

**Observability**
- Request/Response Tracing
- Log Redaction
- Token Usage Monitoring
