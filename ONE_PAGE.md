# CharmingMan — One‑Pager
A privacy‑first, extensible multi‑agent assistant: a local‑first hybrid client that combines a fast Go TUI with optional cloud services for sync, managed embeddings, telephony, and premium LLMs. This one‑pager captures the problem, audience, product scope (MVP), architecture, platform/tech choices, top flows (including step‑by‑step happy paths for the top 3), security & operational constraints, and next steps for product & engineering.

---

## Elevator pitch
CharmingMan is a modular multi‑agent assistant that runs locally (single‑binary TUI) with optional cloud integrations for sync, managed vector DBs, telephony and premium LLMs. It enables power users to do private RAG, transcription, agent orchestration, and tool-enabled automation while giving teams optional managed features (sync, backups, analytics).

---

## Problem statement
Power users and teams need an extensible assistant that can safely process sensitive documents, run agentic workflows, and integrate external tools — all while preserving privacy and offering a clear upgrade path to managed cloud features (sync, telephony, premium LLMs). Current solutions either force cloud-only processing (privacy risk) or are too heavy to run locally and lack multi‑agent orchestration and tool safety.

---

## Target audience & ideal customer
- Primary: Developer and data‑engineering power users who need local-sensitive processing, reproducible automations, and tools integration.
- Secondary: Security/privacy engineers and knowledge workers who must handle sensitive documents (legal, compliance, research).
- Ideal customer profile: Small engineering teams or privacy-conscious freelancers who want a local-first assistant with optional cloud services for sync and advanced features; willing to manage API keys for cloud integrations.

---

## Core value proposition
- Local-first privacy: sensitive data and core processing can remain on the client.
- Modular, agentic orchestration: multiple agents with structured tool calls and stateful chatrooms.
- Extensible tooling with safe defaults: minimal safe toolset for the MVP with explicit confirmations and per-agent allowlists.
- Hybrid path to managed features: optional cloud sync, managed vectors (Pinecone), and premium LLMs (OpenAI) when users opt in.

---

## MVP scope (what we will build first)
Primary goals
- Deliver a productive, local‑first TUI client (Go / Bubble Tea) that can run on desktop.
- Implement multi‑agent orchestration, structured function/tool calling, and stateful chatrooms.
- Provide RAG over uploaded documents (cloud embeddings + managed vector DB by default), plus a documented local-fallback plan.
- Offer provider credential management and cloud sync (optional) backed by a REST sync server.
- Include Whisper transcription hybrid (local transcription own servers) and cloud TTS/telephony integrations deferred to near-term iterations.

Top flows (prioritized)
1. Start conversation (core chat)
2. Upload PDF/doc & RAG QA
3. New agent creation (wizard)
4. Provider credential management (OpenAI / HF / Pinecone)
5. Sync / backup (REST delta sync to cloud)

Step‑by‑step happy paths will be detailed for top 3 (see Flows section).

---

## Key features (MVP)
- Terminal TUI (Go, Bubble Tea + Bubbles + Lipgloss) as primary UI with a small web admin/portal for account/sync settings.
- Local agent manager + in‑process agent orchestration; optional MCP WebSocket server for multi‑client topology.
- LLM providers: OpenAI + Hugging Face Inference (Llama 2 / Mistral hosted models).
- Embeddings & vector DB: OpenAI embeddings + Pinecone managed vector DB (default). Local HNSW + SQLite fallback documented for offline‑only users.
- Document ingestion: PDF & text upload, PDF text extraction, chunking, embedding, index/persist.
- Whisper transcription: local transcription servers (privacy) with cloud TTS (ElevenLabs) and Twilio for telephony optional.
- Tooling: minimal safe‑toolset — HTTP fetcher, web search, file reader + limited shell (explicit user confirmation), and safe external HTTP proxies/whitelists.
- Agent & state: state machine for chatrooms & agents, structured message envelopes (metadata, tool calls, vector refs).
- Persistence: SQLite local DB (SQLDelight recommended) + local artifacts (attachments, models) directory.
- Secrets: OS native keyring for API keys; encrypted file fallback.
- Concurrency: fixed-size worker pool (pre‑fork) one worker per CPU core (min 2, cap 16), shared queue (length 500), 429 when full.

---

## Tech stack (MVP)
- Client: Go TUI — Bubble Tea, Bubbles, Lipgloss, Huh (wizard), Glow (markdown rendering)
- AI orchestration: Fantasy framework (LLM coordination, multi‑agent) + Catwalk (model DB)
- Backend services:
  - Sync server: small REST service for delta sync (SQLite server‑side), user accounts (email/SSO), and optional encrypted storage.
  - MCP WebSocket server for agent messaging/presence (TLS).
- Providers: OpenAI (LLMs & embeddings), Hugging Face Inference, Pinecone, Twilio (deferred), ElevenLabs (deferred).
- Persistence: SQLite (local client DB via SQLDelight), local FS for artifacts/models.
- Packaging: single static Go binary for desktop; optional installers for platform dependencies.

---

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
MVP (0–3 months)
- TUI client with agent creation wizard, core chat, PDF upload & RAG using OpenAI embeddings + Pinecone, provider credential management, local SQLite persistence, OS keyring integration, worker pool.
- Sync REST backend (basic delta sync), optional MCP server (relay only).
- Minimal safe toolset and per-agent allowlist.

Near term (3–6 months)
- HF inference fallback integration.
- Local‑fallback documentation and prototype for HNSW on-disk index.
- Basic TTS (cloud) and Twilio outbound call integration (optional).
- Conflict UI for agent & document merges.

Mid term (6–12 months)
- WASM plugin model for safe tool extensions.
- Mobile/web ports (Flutter) or a web UI backed by the same backend APIs.
- Enterprise features: team sync, audit logs, and compliance tooling.

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
- Embeddings & vector DB: OpenAI embeddings + Pinecone (managed) by default; local HNSW fallback documented.
- Transcription/TTS/telephony: Hybrid — local transcription (self‑run servers) + cloud TTS (ElevenLabs) and Twilio for telephony (deferred).
- Secrets: OS keyring; encrypted fallback.
- Persistence: SQLite local DB (SQLDelight recommended) + artifacts folder.
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
