# CharmingMan AI Gateway API Documentation

The CharmingMan AI Gateway provides a unified REST interface for multiple Large Language Model (LLM) providers, Retrieval-Augmented Generation (RAG), and multimodal capabilities.

## Base URL
`http://localhost:8090/api/v1`

## Authentication

Authentication is handled via a custom API key header. If the gateway is configured with a `GATEWAY_API_KEY`, all requests must include:

| Header | Description |
|--------|-------------|
| `X-Charming-Key` | Your AI Gateway API key |

---

## 1. Chat
`POST /chat`

Unified chat endpoint supporting multiple providers, persistent history, and optional RAG.

### Request Body
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `provider` | string | Yes | LLM provider (e.g., `openai`, `anthropic`, `ollama`, `llamacpp`) |
| `model` | string | Yes | Model identifier (e.g., `gpt-4o`, `claude-3-5-sonnet-20240620`) |
| `prompt` | string | Yes | The user's input message |
| `use_rag` | boolean | No | Enable semantic search over indexed documents |
| `room_id` | string | No | UUID of the chat room for history persistence |
| `agent_id` | string | No | ID of a pre-configured agent to use |

### Response Schema
```json
{
  "response": "The generated text response from the model",
  "usage": {
    "prompt_tokens": 120,
    "completion_tokens": 45,
    "total_tokens": 165
  },
  "sources": [
    {
      "document_id": "uuid-123",
      "content": "Relevant document snippet...",
      "score": 0.89
    }
  ],
  "error": null
}
```

---

## 2. Usage
`GET /usage`

Retrieve token usage logs and aggregate statistics.

### Response Schema
```json
{
  "logs": [
    {
      "id": 1,
      "provider": "openai",
      "model": "gpt-4o",
      "total_tokens": 70,
      "latency_ms": 450,
      "timestamp": "2023-10-27T10:00:00Z"
    }
  ],
  "stats": {
    "total_tokens": 15400,
    "total_cost": 0.45,
    "total_requests": 210
  }
}
```

---

## 3. Agents
`GET /agents` | `POST /agents`

Manage agent personas and model configurations.

---

## 4. Knowledge (RAG)
`POST /documents` | `GET /search`

Index documents and perform semantic search.

---

## 5. Multimodal

### Transcribe
`POST /transcribe`

Accepts `multipart/form-data` with an audio file (under the `file` key). Returns transcribed text.

### Speech
`POST /speech`

Accepts JSON `{ "text": "..." }`. Returns an MP3 audio stream.

---

## 6. Tooling
`GET /tools`

Returns a list of all active Model Context Protocol (MCP) tools discovered from local servers.

```json
[
  {
    "name": "read_file",
    "description": "Reads the content of a local file."
  },
  {
    "name": "list_directory",
    "description": "Lists files in a specified directory."
  }
]
```
