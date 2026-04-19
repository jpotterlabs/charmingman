an extensible multi‑agent assistant: a local‑first hybrid client that combines a fast Go TUI with optional cloud services for sync, managed embeddings, telephony, and premium LLMs.
Build a modular multi‑agent assistant that runs locally (single‑binary TUI) with cloud integrations. It enables power users to do RAG, transcription, agent orchestration, TTS, STT, and telephony. Key features include agent creation, room creation, and multi-agent orchestration with human orchestration using @mentions. Different agents within the same room are aware and can see other agents' responses. The human coordinator will call on which agent speaks next using @mentions.

Process  documents, run agentic workflows, and integrate external tools — have all core features powered by either local or remote openai compatible keys

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
Whisper transcription: local transcription servers (privacy) with cloud TTS (ElevenLabs) and Twilio for telephony optional.
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
2. Upload PDF/doc & RAG QA
3. New agent creation (wizard)
4. Provider credential management (OpenAI / HF / Pinecone)
5. Sync / backup (REST delta sync to cloud)
## Architecture overview (high level)
- Client (TUI)
  - Presentation: Bubble Tea UI
  - Local Agent Manager: in‑process orchestrator and state machine
  - Model Manager: handles local models, downloads, integrity checks
  - Persistence: SQLite DB + artifacts folder
  - Secrets: OS keyring adapter
- AI Gateway (Fantasy + provider adapters)
  - Abstract LLM / embeddings / tools interface
  - Provider adapters: OpenAI, HF Inference, Pinecone client
  - Function/tool invocation coordinator (structured JSON schema)
- Cloud Backend
  - Sync REST API (delta sync, conflict handling)
  - MCP WebSocket server (message coordination)
  - Optional (encrypted) server storage when users opt into sync
- External Providers: OpenAI, Hugging Face Inference, Pinecone, Twilio, ElevenLabs

API boundaries and contracts
- Small, well‑typed JSON schema / protobufs between components:
  - Message envelope: id, timestamp, agent_id, role, content[], tools[], vector_refs[], metadata
  - Tool call spec: tool_id, args (typed), expected outputs, dry-run flag, approval token
  - Agent manifest: name, persona, tool allowlist, memory policy, sync flag

---

## Flows (Top 5) — happy path details

1) Start conversation — step‑by‑step (happy path)
- User opens CharmingMan (TUI) and selects "New Conversation" or starts from agent list.
- Client loads selected agent manifest (persona, tool allowlist) and local context (pinned docs).
- User types message; TUI displays streaming response from LLM adapter.
- Fantasy gateway formats message envelope, performs reasoning/chain steps and tool calls as required.
- If a tool call requires side-effects (e.g., HTTP POST, limited shell), the client shows a confirmation dialog. If approved, the tool runs and outputs are returned to the agent.
- Conversation and action metadata persist to SQLite; background job queue handles any long-running tasks (indexing, embedding) and shows status in UI.
- If sync is enabled, the client enqueues a delta and pushes to sync server when connectivity permits.

2) Upload PDF/doc & RAG QA — step‑by‑step (happy path)
- User selects "Upload Document" → selects file(s).
- Client extracts text (PDF extractor / OCR hook), chunks content (configurable chunk size + overlap) and creates local artifact entries.
- Default flow (MVP): client uploads raw text/chunks or embeddings to managed services (if user opted in):
  - If remote embeddings chosen: client requests OpenAI embeddings for chunks, then upserts vectors into Pinecone, storing vector ids in local SQLite metadata.
  - If user has offline mode enabled: client optionally computes local embeddings (future feature) and persists HNSW index locally.
- After indexing is complete, user opens a chat and asks document‑specific query.
- The RAG pipeline:
  - Query → embed (OpenAI) → Pinecone query → retrieve top N chunks → assemble context + prompt → LLM call (OpenAI / HF fallback).
  - LLM generates response; client displays source citations and provides a "View Source" command to open original chunk.
- All document metadata (source, extraction timestamps, embeddings status) saved in local DB.

3) New agent creation (wizard) — step‑by‑step (happy path)
- User chooses "Create Agent" → Huh wizard guides:
  - Step 1: Name & description
  - Step 2: Persona & instructions (system prompt editor)
  - Step 3: Tool allowlist (HTTP fetch, file read, limited shell — default: HTTP & file; shell disabled)
  - Step 4: Memory & context policy (ephemeral, pinned docs, long-term memory toggle)
  - Step 5: Sync settings (local-only, sync-enabled)
  - Step 6: Review & create
- Agent manifest saved locally; user can start a conversation with new agent immediately.
- First time a tool is required and not yet authorized, the client prompts for per‑agent allowlist or per-command confirmation.

4) Provider credential management (high level)
- From settings, user adds provider keys (OpenAI, HF token, Pinecone API key).
- Keys are stored in OS keyring; a config pointer stored in SQLite referencing keyring entry.
- For headless/machine usage, support machine API keys and allow export/import via encrypted token file.

5) Sync / backup (high level)
- Optional REST sync: client sends deltas using per-row lastUpdated timestamps and a signed token from user account.
- Conflict policy default: last‑write‑wins; for agent manifests & documents show a simple conflict UI to resolve differences before accepting server version.
- Server persists encrypted content only when user explicitly enables cloud backup for an item; otherwise server stores metadata only.

---

## Offline vs cloud boundaries
- Fully offline features (if local model & local index present): conversation with local model, RAG over local documents, local transcription (if user runs local transcription service), local tools (file read, limited shell).
- Cloud default for MVP: embeddings (OpenAI), vector DB (Pinecone), LLMs (OpenAI primary), HF inference fallback. Local-fallback plan (HNSW + on-disk embeddings, ggml/llama.cpp execution) documented as low-priority feature.
- Sensitive data rule: plaintext chats, uploaded documents, and audio never leave the client unless user explicitly opts in per-item or enables sync.

---

## Security, secrets & compliance
- Secrets: OS-native keyring (macOS Keychain, Windows Credential Manager, Linux Secret Service). Fallback: encrypted config file (AES derived from passphrase).
- Tool safety:
  - Minimal safe‑toolset for MVP.
  - Per‑agent allowlist + first-use per-command confirmation for shell and side-effect tools.
  - External HTTP requests go through an optional proxy and domain whitelist.
  - Dry-run toggle for side-effects.
- Telemetry & privacy:
  - Opt‑in telemetry only. If opted in, collect anonymized crashes, feature counts, API error rates, hashed environment id; 90‑day retention; user can view/export/delete telemetry data.
- MCP/WebSocket: by default messages relayed only, not persisted. If user opts into cloud sync, encrypted content is stored server‑side (encrypted at rest) and audit logs are available if enabled.
- Compliance: No enterprise compliance (GDPR/HIPAA/SOC2) in MVP — but provide export/delete primitives to enable later certification.

---

## Data model & persistence
Local persisted items:
- Chats & messages (id, agent_id, role, content[], metadata)
- Agent manifests & configs
- Attachment refs & artifact metadata
- Indexing metadata (chunks, vector ids)
- Job queue (background tasks)
- Local model metadata
Storage:
- SQLite via SQLDelight for structured data
- Artifacts directory for files & downloaded models
- Optional encrypted blob storage when item is marked for cloud backup

Sync semantics:
- REST delta sync with per-row lastUpdated timestamps; last‑write‑wins by default.
- Conflict UI for agent manifests and documents.

---

## Concurrency & worker model
- Fixed-size pre‑fork worker pool: one worker per CPU core (min 2, cap 16).
- Each worker runs an event loop; round‑robin load balancing across the pool.
- Shared incoming queue length default = 500; reject with 429 when full.
- Rate limiters per provider (configurable) and persistent job queue (SQLite) for long-running tasks (indexing, downloads, inference).

---

## Packaging & distribution
- MVP client: single static Go binary for macOS and Linux (CI cross-compile for Windows soon).
- Provide installers or scripts for model dependencies (if local models supported).
- Model download strategy: hybrid — auto-download small models (<200 MB) with user consent for large models (≥200 MB). Show size/time estimate and Wi‑Fi only scheduling; integrity checks and optional auto‑update.

---

## Testing & CI
- Unit tests for provider adapters (Fantasy mocks), agent state machine, persistence layer.
- Integration tests for cloud provider flows using sandbox keys (OpenAI/HF/Pinecone test accounts).
- E2E tests with terminal emulation for TUI (recorded sessions) and job queue/integration tests for background tasks.
- Nightly builds and a monthly stable release cadence.

---

## Roadmap & near‑term milestones
Phase 1 - 3: Completed ✅
- TUI client with agent creation wizard, core chat, PDF upload & RAG using OpenAI embeddings + Pinecone, provider credential management, local SQLite persistence, OS keyring integration, worker pool.
- Initial Window Management ("The Stage") and draggable TUI components.

Near term (Phase 4): 🚀
- Multi-Agent Swarms & Tools: Agent-to-Agent communication.
- MCP (Model Context Protocol) integration for local tool calling.
- Thinking Drawer for visualizing reasoning chains.

Mid term (Phase 5):
- HF inference fallback integration.
- Local‑fallback documentation and prototype for HNSW on-disk index.
- Basic TTS (cloud) and Twilio outbound call integration (optional).
- Conflict UI for agent & document merges.

---

## Business & licensing
- Monetization: freemium — free local client; paid subscription for cloud sync, managed vector DB usage, telephony, team features.
- Accounts: optional; local-only use allowed. Account creation (email+password + optional SSO) only required for cloud features.
- License: MIT for codebase. Explicit policy: third‑party model weights must respect upstream licenses; do not bundle proprietary model weights in repo.

---

## Decisions log / chosen defaults (summary)
- Distribution: Hybrid — local‑first client with optional cloud backend for sync/managed services.
- Primary UI for MVP: Terminal TUI in Go (Bubble Tea) + small web admin portal.
- LLMs: OpenAI + Hugging Face Inference (HF) as fallback.
- Embeddings & vector DB: OpenAI embeddings + Pinecone (managed) by default.
- Transcription/TTS/telephony: Hybrid — local transcription (self‑run servers) + cloud TTS (ElevenLabs) and Twilio for telephony (deferred).
- Secrets: OS keyring; encrypted fallback.
- Persistence: SQLite local DB + artifacts folder.
- Concurrency: pre‑fork worker pool (1 per CPU core, min 2, cap 16), queue 500.
- Telemetry: opt‑in only, anonymized, 90‑day retention.
- License: MIT.

---

## Open questions & engineering / product action items
1. Confirm primary client for MVP — proceed with TUI + small web portal (recommended default), or pivot to Flutter mobile/web first? (Default: TUI)
2. Confirm Pinecone as the initial managed vector DB (or pick Weaviate/Redis). (Default: Pinecone)
3. Define the exact JSON schema for message envelopes and tool call contracts (engineering to propose v1).
4. Confirm size threshold for auto-download model policy (suggested default: 200 MB).
5. Define production plan for MCP server persistence policy (relay-only vs. persisted encrypted content) and retention defaults.
6. Create a minimal CLI admin plan for headless / CI installs for keyring fallback.
7. Prepare legal guidance for model download assistant (link to model licenses & user consent flow).

---

## Next steps (recommended)
- Product: approve this one‑pager and prioritize the feature list (MVP vs v1).
- Engineering:
  - Draft API schemas (message envelope, tool call, agent manifest).
  - Prototype agent manager & Fantasy adapter to OpenAI + HF.
  - Implement local SQLite persistence schemas (chat, agent_manifest, artifacts, jobs).
  - Build credential storage using OS keyring and an encrypted fallback for headless installs.
  - Implement initial TUI screens: agent list, new agent wizard, conversation screen, doc upload/import.
- Security/Legal: prepare telemetry privacy policy, terms covering model downloads and redistribution, and account/SSO flows.
- Ops: prepare a small sync server prototype (REST delta sync) and a minimal MCP WebSocket relay.

---

If you want, I can now:
- Produce the API schemas (message envelope, tool descriptor, agent manifest) as concrete JSON/protobuf drafts.
- Generate screen-by-screen TUI wireframes & UI copy for the top 3 flows.
- Produce an engineering kickoff checklist with prioritized tickets for a 4–6 week sprint.

Which would you like next?