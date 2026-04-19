an extensible multi‑agent assistant: a local‑first hybrid client that combines a fast Go TUI with optional cloud services for sync, managed embeddings, telephony, and premium LLMs.
Build a modular multi‑agent assistant that runs locally (single‑binary TUI) with cloud integrations. It enables power users to do RAG, transcription, agent orchestration, TTS, STT, and telephony. Key features include agent creation, room creation, and multi-agent orchestration with human orchestration using @mentions. Different agents within the same room are aware and can see other agents' responses. The human coordinator will call on which agent speaks next using @mentions.

Process documents, run agentic workflows, and integrate external tools — have all core features powered by either local or remote openai compatible keys

Modular, agentic orchestration: multiple agents with structured tool calls and stateful chatrooms.
Deliver a productive, local‑first TUI client (Go / Bubble Tea) that can run on desktop.
Implement multi‑agent orchestration, structured function/tool calling, and stateful chatrooms.
Provide RAG over uploaded documents 
Offer provider credential management
Include Whisper transcription hybrid (like all other services)


Terminal TUI (Go, Bubble Tea + Bubbles + Lipgloss) as primary UI with a small web admin/portal for account/sync settings.
Local agent manager + in‑process agent orchestration;
LLM providers: local- ollama, vllm, llama.cpp; remote- openrouter, openai, anthropic
Embeddings & vector DB: OpenAI text-embedding-3-small for embeddings + Pinecone for managed storage and ANN search — configured with OPENAI_API_KEY and PINECONE_API_KEY.
Document ingestion: PDF & text upload, PDF text extraction, chunking, embedding, index/persist.
Whisper transcription: Integrated Whisper STT via AI Gateway (Phase 5).
Agent & state: state machine for chatrooms & agents, structured message envelopes (metadata, tool calls, vector refs).
Separate Concerns
Frontend TUI is fully agnostic
Backend API-Gateway is fully agnostic, providing a unified api abstracting to hundreds of llms
Client: Go TUI — Bubble Tea, Bubbles, Lipgloss, Huh (wizard), Glow (markdown rendering)
AI orchestration: Fantasy framework (LLM coordination, multi‑agent) + Catwalk (model DB)
Providers: OpenAI compatible remote or local 
Local Agent Manager: in‑process orchestrator and state machine
Model Manager: handles local models, downloads, integrity checks




Top flows (prioritized)
1. Start conversation (core chat)
2. Voice-to-Text interaction (new in Phase 5)
3. Upload PDF/doc & RAG QA
4. New agent creation (wizard)
5. Provider credential management (OpenAI / HF / Pinecone)
## Architecture overview (high level)
- Client (TUI)
  - Presentation: Bubble Tea UI
  - Voice Suite: System-level recording via `sox` (v5)
  - Local Agent Manager: in‑process orchestrator and state machine
  - Persistence: SQLite DB + artifacts folder
- AI Gateway (Fantasy + provider adapters)
  - Abstract LLM / embeddings / tools interface
  - Transcription Handler: Whisper STT (v5)
  - Provider adapters: OpenAI, HF Inference, Pinecone client
- Cloud Backend
  - Sync REST API (delta sync, conflict handling)
  - MCP WebSocket server (message coordination)
- External Providers: OpenAI, Hugging Face Inference, Pinecone, Twilio, ElevenLabs

API boundaries and contracts
- Small, well‑typed JSON schema / protobufs between components:
  - Message envelope: id, timestamp, agent_id, role, content[], tools[], vector_refs[], metadata
  - Transcription request: multipart/form-data with audio file (v5)

---

## Roadmap & near‑term milestones
Phase 1 - 5: Completed ✅
- TUI client with agent creation wizard, core chat, PDF upload & RAG.
- Infinity Canvas with panning, zooming, and spatial organization (Phase 4).
- **Voice & Multimodal**: Whisper STT integration and 'v' key recording (Phase 5).
- **Stability**: Bounded chat history (10 msgs), log redaction, and coordinate scaling fixes (Phase 5).

Near term (Phase 6): 🚀
- MCP (Model Context Protocol) integration for local tool calling.
- Text-to-Speech (TTS) for agent responses.
- Thinking Drawer for visualizing reasoning chains.

---

## Decisions log / chosen defaults (summary)
- Distribution: Hybrid — local‑first client with optional cloud backend.
- Primary UI: Terminal TUI in Go (Bubble Tea).
- LLMs: OpenAI + Hugging Face Inference (HF) as fallback.
- Embeddings & vector DB: OpenAI embeddings + Pinecone (managed) by default.
- Transcription/TTS/telephony: Whisper STT (v5) + ElevenLabs TTS (Planned).
- Secrets: OS keyring; encrypted fallback.
- Persistence: SQLite local DB + artifacts folder.
- License: MIT.
