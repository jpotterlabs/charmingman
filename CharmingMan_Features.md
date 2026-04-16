# Project Blueprint: CharmingMan
**CharmingMan** is a sophisticated, multi-agent ChatTUI built exclusively on the **Charm** ecosystem. It leverages a modular Go-based architecture to provide a seamless bridge between local/remote LLMs and a high-performance terminal interface.

---


- CharmingMan: Technical Architecture & Feature Expansion

    - The **CharmingMan** project is an ambitious use of the Charmbracelet ecosystem. By leveraging Go’s concurrency and the TUI's inherent focus on text, you can build a tool that rivals web-based AI dashboards while remaining lightweight and terminal-native.

    - The application follows a strict separation of concerns, utilizing a layered approach to manage complex agentic workflows and a high-fidelity terminal UI.
---

## The Frontend (TUI Layer)

---

### The CharmingMan "Command Center" Design

1. The Multi-Pane "Tiled" Workspace
Instead of a single column of text, the main dashboard uses a **Flexbox-style layout** powered by `Lipgloss`.

* **The Quad-View:** Support for up to four active chat panes simultaneously. Each pane is an independent `Bubble Tea` model.
* **Active Focus:** The "focused" pane is highlighted with a bold, `thicker` border and a specific color (e.g., Mauve for GPT-4, Green for local Ollama) using `Catwalk` to ensure the colors pop against the terminal background.
* **The "Beam" TUI Implementation:** A dedicated horizontal split where the top 70% shows the multi-model responses side-by-side, and the bottom 30% is a shared input buffer for "scattering" the prompt.



---

2. The Sidebar "Intelligence Hub"
A collapsible left-hand navigation bar (toggled with `ctrl+b`) manages the application's global state.

* **Agent Library:** A vertical list of custom agents (created via the `Huh?` wizard). Icons next to their names indicate the backend (e.g., a "☁️" for Anthropic, a "🏠" for local Llama).
* **Active MCP Servers:** A status indicator showing which Model Context Protocol servers are currently attached (e.g., "Filesystem: Active", "PostgreSQL: Connected").
* **Token Telemetry:** A sparkline-style mini-graph in the bottom corner of the sidebar showing real-time token throughput or cost-per-session.

---

3. The "Stage" & Document Engine
The right-hand side of the dashboard is reserved for **The Stage**, a context-aware viewport for non-chat artifacts.

* **Document Preview:** When you "Chat with a PDF," the document text is rendered in this pane using `Glow`. As the AI references specific pages, the `Glow` viewport auto-scrolls to the relevant section.
* **Code & Artifacts:** If the agent generates code or a "DevSpec," it is promoted to The Stage. This keeps the chat history clean while allowing you to inspect the output in full-screen width.
* **Live Reasoning Log:** A "Thinking" drawer that shows the hidden chain-of-thought, tool-calling logs, and search results without cluttering the primary conversation.

---

4. Interactions & Navigation
The dashboard prioritizes keyboard-driven speed, mimicking the efficiency of an IDE like Neovim.

| Keybinding | Action | Component |
| :--- | :--- | :--- |
| `Tab` | Cycle focus between active chat panes | `Bubble Tea` |
| `Ctrl + N` | Open the `Huh?` Wizard for a new Agent/Chat | `Huh?` |
| `Ctrl + S` | Toggle the "Stage" (Document/Glow view) | `Lipgloss` |
| `/` | Enter "Command Mode" for tool calling or model switching | `Bubbles` |

---

5. Visual Aesthetic
To achieve that professional, high-end feel in a terminal:
* **Adaptive Theme:** Use `Catwalk` to sync the TUI theme with the user's terminal (Dracula, Solarized, Nord) while injecting specific accent colors for different AI providers.
* **Status Lines:** A persistent bottom bar (similar to `vim-airline`) showing the current model, temperature, active "tools" (e.g., Whisper, MCP), and system latency.
* **Glow-Enhanced Response:** AI responses aren't just strings; they are processed through `Glow` to provide syntax highlighting for code and beautiful tables for structured data.



---

### 2. The AI Gateway (Middleware)
This is the most critical layer for modularity. By using **Fantasy**, you are already positioned to handle structured data.

* **Agent Coordination:** Define agents as **Protocols**. An agent isn't just a prompt; it’s a set of capabilities (Tools + Model + State).
* **Provider Translation:** * Create a `Provider` interface in Go. This allows you to swap between `OpenAI`, `Anthropic`, or `Local` (Ollama/Llama.cpp) without changing the Gateway logic.
    * **Streaming Logic:** Ensure your Gateway supports **Server-Sent Events (SSE)** or equivalent streaming patterns so the TUI can update the chat bubbles in real-time as tokens arrive.

---

### 3. The Backend (Document & Multi-Modal Engine)
This layer handles the "heavy lifting" outside of the main UI loop.

#### **RAG & Document Management**
* **Vector Store:** For a local-first TUI, consider a lightweight vector store like **ChromaDB** or even a local **FAISS** index stored as a file.
* **PDF Processing:** Use a Go-native library like `unidoc/unipdf` to extract text. The Gateway then chunks this text and sends it to the embedding model.

#### **Audio & Communication (The "Omni" Features)**
To handle Whisper (STT), TTS, and Phone calls without locking the UI, these must run in **background Goroutines**.

| Feature | Implementation Strategy |
| :--- | :--- |
| **Whisper (STT)** | Capture audio via `portaudio` or `ffmpeg`, stream to a Whisper endpoint, and send a `WhisperCompleteMsg` back to Bubble Tea to populate the input field. |
| **TTS** | Stream response text to a TTS provider (like ElevenLabs or OpenAI) and play audio through a local buffer while the text renders in the TUI. |
| **Phone Calls** | Integrate **Twilio** or **SignalWire** APIs. The AI Gateway treats the phone call as a specialized "session" where input comes from a webhook rather than the keyboard. |
| **Tool Calling** | Use **JSON Schema** definitions. When the model requests a tool (e.g., `get_weather`), the Gateway executes the local Go function and feeds the result back to the LLM before the UI even sees the final answer. |

---

### Proposed State Machine Flow


1.  **Idle State:** User views the dashboard (New Chat / Load Session).
2.  **Configuration State (Huh):** User defines agent parameters.
3.  **Active Chat State (Bubble Tea):**
    * **Sub-state: Awaiting Input**
    * **Sub-state: Processing (Spinner component)**
    * **Sub-state: Streaming Response**
4.  **Document State (Glow):** User reviews a generated `DevSpec` or `CharterDoc`.
5.  **Action State:** A tool is triggered (e.g., making a phone call or searching a local PDF).

---

### Technical Tip: The "Msg" Pattern
In Bubble Tea, everything is a `Cmd`. To keep the backend clean, your Backend services should return functions that wrap the result in a custom message type:

```go
func (b *Backend) TranscribeAudio() tea.Cmd {
    return func() tea.Msg {
        text := b.WhisperClient.Run()
        return TranscriptionResultMsg(text)
    }
}
```
This ensures your TUI stays responsive (60fps) while the AI Gateway or Document Engine is doing complex processing in the background.




---

### 3. The AI Gateway (Middleware)
This layer serves as the "brain," decoupling the UI from the underlying intelligence providers.

* **Agentic Framework (`Fantasy`):**
    * Manages the lifecycle of multi-agent coordination.
    * Enforces **Structured Responses** by mapping AI outputs to Go structs for predictable UI updates.
* **Provider Translation Layer:**
    * A unified interface to route requests to local (Ollama/LocalAI) or remote (OpenAI, Anthropic, Hugging Face) endpoints.
* **Model Management:**
    * Integration with **Hugging Face** for direct model downloading and local storage.
    * Support for user-provided GGUF or Safetensors files.

---

### 4. The Backend & Document Engine
The heavy-lifting layer responsible for persistence, multimodal capabilities, and external connectivity.

#### **State & Memory Management**
* **Stateful Chatrooms:** Conversations are stored server-side (or in a local DB) so agents maintain context without requiring the full history to be sent in every payload.
* **Prompt Processing:** Pre-processing of user input for reasoning chains and multi-step tool use.

#### **Multimodal & Communication Services**
* **Document Intelligence:** RAG (Retrieval-Augmented Generation) pipeline supporting PDF/Document uploads and text embedding.
* **Voice Suite:** * **Whisper Integration:** Local/API transcription of audio files.
    * **TTS (Text-to-Speech):** Audio synthesis for agent responses.
    * **Telephony:** Integration for making outbound phone calls via agentic triggers.

#### **Agentic Connectivity**
* **MCP & Protocols:** Support for the **Model Context Protocol (MCP)** to allow agents to interact with local files and databases.
* **Websockets/A2A:** Enables real-time Agent-to-Agent (A2A) communication and external tool calling.

---

### 5. Implementation Roadmap

| Phase | Focus | Key Tools |
| :--- | :--- | :--- |
| **I: Foundations** | Project scaffolding and basic Bubble Tea loop. | Bubble Tea, Lipgloss |
| **II: The Wizard** | Agent configuration and provider setup. | Huh?, Fantasy |
| **III: Intelligence** | Provider translation and structured responses. | Fantasy, Catwalk |
| **IV: Knowledge** | PDF uploads, Vector storage, and RAG. | Go-Document-Parsers |
| **V: Expansion** | TTS, Whisper, and MCP Tool Calling. | Websockets, Whisper API |





To add a roleplay dimension to **CharmingMan**, you can leverage the **Fantasy** framework's structured response capability to "bake in" personality traits during the agent creation phase. 

By treating "Personality" as a first-class data object rather than just a text block, your agents will behave with much higher consistency in the TUI.

---

### 1. The Personality Schema
Instead of a generic system prompt, define a `Personality` struct. This allows the AI Gateway to "tune" the model’s behavior programmatically.

```go
type Personality struct {
    Name       string   `json:"name"`
    Archetype  string   `json:"archetype"` // e.g., "Stoic Mentor", "Sarcastic Tech Support"
    Traits     []string `json:"traits"`    // e.g., ["verbose", "witty", "cautious"]
    Speech     string   `json:"speech"`    // e.g., "Pirate slang", "Formal Victorian"
    CoreValues []string `json:"values"`    // Guides how the agent makes decisions
}
```

### 2. Updating the "Huh" Wizard
In your agent creation phase, add a specific **Persona Workshop**. Using `huh.Select` and `huh.MultiSelect`, you can let users build a personality without typing a wall of text.

* **Archetype Selection:** Use a list of presets (e.g., The Rebel, The Sage, The Assistant).
* **Tone Sliders:** While `huh` is largely text-based, you can use a `huh.Select` to set a "Temperament" (e.g., Aggressive <—> Passive).
* **Voice/Speech Toggles:** Use `huh.Confirm` to toggle specific quirks like "Uses Emojis" or "Never uses capital letters."

### 3. Implementing the "Persona Engine"
In your **AI Gateway (Middleware)**, you need a function that compiles these traits into a "Super-System Prompt." 

> **Tip:** Do not just append the traits. Use a template to give them weight.

```markdown
# Roleplay Protocol: {{.Name}}
You are currently roleplaying as {{.Name}}, a {{.Archetype}}.
## Constraints:
- Tone: {{.Traits}}
- Linguistic Style: {{.Speech}}
- Behavioral Bounds: {{.CoreValues}}

## Interaction Instructions:
Never break character. If the user asks about your AI nature, respond in the voice of {{.Name}}.
```

### 4. TUI Visual Cues (Lipgloss & Bubbles)
To make roleplaying feel immersive in a terminal, the interface should visually reflect who you are talking to.

* **Dynamic Borders:** Use `Lipgloss` to change the border color of the chat bubble based on the agent's personality. 
    * *Example:* A "Hot-Headed" agent gets a `Red` border; a "Calm" agent gets `ForestGreen`.
* **Avatar Rendering:** Since you are in the terminal, use **ASCII Art** or **ANSI icons** in the chat header. You can store a small library of ASCII avatars that match the `Archetype` selected in the Wizard.
* **Status Bars:** Use a bottom status bar to show the agent's current "Mood" or "State" if your state machine tracks emotional shifts during the conversation.

### 5. Advanced: Memory & Continuity
For deep roleplay, the agent needs to remember past interactions.
* **The "Vibe" Buffer:** In your **Backend (Document Engine)**, store a "Memory Summary" in your RAG system. Instead of just searching for facts, search for "Relational Data"—how the agent feels about the user based on previous chats.
* **Fantasy Integration:** Use Fantasy's `Structured Response` to have the model output two things: the **In-Character Response** and a hidden **Emotional State Update** that you save to the database.

---

### Example Archetype: The Roman Senator
Since you appreciate historical accuracy, you could have a preset for a **Roman Senator (90 BCE)**:
* **Traits:** ["Authoritative", "Traditionalist", "Oratorical"]
* **Speech:** "Latinate sentence structures, refers to the Mos Maiorum frequently."
* **TUI Style:** Purple (`latus clavus`) borders and a marble-white text color.


To elevate **CharmingMan** from a standard chat application to a true roleplay environment, you can implement a **World Book** (often called a "Lorebook" or "Codex"). This acts as a secondary RAG (Retrieval-Augmented Generation) layer specifically for setting-specific facts, history, and rules that the agent should know but shouldn't necessarily keep in its immediate "brain" (context window) at all times.

---

### 1. The Lore Engine Architecture
The World Book should sit within your **Backend (Document Engine)**, but unlike the standard PDF upload, it uses **Keyword-Triggered Injection**.

* **Lore Entries:** A collection of snippets (Markdown files or JSON) containing a `Primary Key` (e.g., "The Republic"), `Aliases` (e.g., "Rome", "SPQR"), and the `Content`.
* **The Scanner:** Before sending the user’s message to the AI Gateway, the Backend scans the text for these keys.
* **Dynamic Injection:** If a match is found, the relevant lore snippet is injected into the System Prompt as "Temporary Context."

---

### 2. World Book Management (The TUI Interface)
Using the **Charm** ecosystem, you can make managing lore as satisfying as the roleplay itself:

* **The Codex View (`Bubbles.List`):** A searchable sidebar where users can scroll through all "Known Lore."
* **The Entry Editor (`Huh`):** A specialized form to create new entries.
    * **Fields:** Category (Place, Person, Event, Item), Description, and "Importance Level" (decides if the info is always in context or only when mentioned).
* **Markdown Preview (`Glow`):** Since World Books often involve long descriptions, use `Glow` to render the history of a city or the details of a specific era (like the Social War of 90 BCE) with proper terminal formatting.

---

### 3. Implementation Logic
Here is how the data flow looks when a user mentions something from the World Book:



1.  **Input:** User types "What do you think of the new laws passed in the Forum?"
2.  **Detection:** The Backend identifies **"Forum"** and **"laws"** as keys in the World Book.
3.  **Fetch:** It pulls the entry: *"The Forum Romanum is the heart of political life... current laws under debate include the Lex Julia."*
4.  **Augment:** The AI Gateway receives the message with a hidden prefix: `[Lore Reference: The Forum is currently debating the Lex Julia]`.
5.  **Output:** The agent responds with specific, lore-accurate details, maintaining the "illusion" of a living world.

---

### 4. Advanced Features: "The Fog of War"
To make the TUI feel like a game or a deep simulation, you can add state-based lore:

* **Hidden Entries:** Some lore entries are "Locked." They only become "Known" (visible in the TUI) once the agent mentions them in conversation.
* **Consistency Checker:** Use **Fantasy's Structured Response** to let the agent *propose* new lore. If you are roleplaying and the agent invents a name for a tavern, CharmingMan can ask: *"Add 'The Rusty Tankard' to the World Book?"*
* **Chronological Lore:** If your world-building involves a timeline, the World Book can filter entries based on a "Current Date" variable in your state machine.

---

### 5. Technical Structure
In your Go backend, you can represent the World Book as a simple key-value store or a local SQLite database for faster lookups:

```go
type LoreEntry struct {
    ID          string   `json:"id"`
    Keywords    []string `json:"keywords"`
    Body        string   `json:"body"`
    IsDiscovered bool     `json:"is_discovered"`
}
```

This allows you to keep the **AI Gateway** clean while ensuring the **CharmingMan** agents are the most well-informed "personalities" in the terminal.


help me flesh out this idea CharmingMan, a Multi-Agent ChatTUI
Charm ecosystem. modular Go application with a clean split between the terminal interface and the backend logic.
Technical Architecture
To ensure scalability and clean code, we will divide the application into three primary layers:
1. The Frontend (TUI Layer)
Huh & Bubbles: Used for the agent creation Wizard" phase. Huh handles walking through user creating an agent and the various settings
Bubble Tea: The main orchestration engine. It will manage the main chatUI interface (New Chat, New Agent, Upload PDF)
Lipgloss: Used to define the visual style, borders, and padding for a professional TUI look.
Glow & Catwalk: Glow will render the generated Markdown docs (CharterDoc, DevSpec, etc.) directly in the terminal for user approval. Catwalk ensures the TUI remains compatible with various terminal themes and LLM-derived color palettes.
2. The AI Gateway (Middleware)
Fantasy Framework: Acts as the "brain." It will handle providers, endpoint mapping, agent coordination, etc. Structured Response capability will define models for the AI to follow. 
Provider Translation: A distinct service that takes the internal request and routes it to the correct endpoint (Local, OpenAI, Anthropic) using keys from .env.local.
3. The Backend (Document Engine)
State Machine: Tracks which state the user is in
RAG and Question Ansering
support document and PDF upload, chat with PDF, and Question Answering
Support transcrobe via whispr
Support TTS
Support outbound phone call making
Support tool calling

---

# gemini response

## CharmingMan: Technical Architecture & Feature Expansion

The **CharmingMan** project is an ambitious use of the Charmbracelet ecosystem. By leveraging Go’s concurrency and the TUI's inherent focus on text, you can build a tool that rivals web-based AI dashboards while remaining lightweight and terminal-native.

Here is a breakdown of how to refine the architecture and integrate the advanced features you mentioned.

---

### 1. The Frontend (TUI Layer)
To maintain a "clean split," the TUI should act as a **State Observer**. It shouldn't know *how* an agent is created, only how to display the form and send the results back to the Gateway.

* **The Wizard (Huh):** Use `huh.Form` for the initial agent configuration (Name, Role, Model selection, System Prompt). This keeps the setup phase distinct from the chat experience.
* **The Orchestrator (Bubble Tea):** * Implement a **View Switcher** pattern. Your `Update` function should delegate messages to different "screens" (e.g., `ChatView`, `DocumentView`, `SettingsView`).
    * **Viewport Optimization:** For the chat interface, use the `viewport` bubble to handle scrolling through long histories, ensuring smooth performance even with large context windows.
* **Visual Feedback (Lipgloss & Catwalk):**
    * Use `catwalk` to test your layouts across light/dark modes. Since LLM responses often contain syntax highlighting, ensure your `Lipgloss` styles use **ANSI colors** or adaptive hex codes to avoid unreadable text on different terminal themes.
* **Onboarding & Configuration (`Huh?`):**
    * A guided **Agent Creation Wizard** to define personas, select LLM providers, and configure API keys.
    * Dynamic forms for system prompt engineering and tool-calling permissions.
* **Orchestration (`Bubble Tea`):**
    * The central The `Model-View-Update` (MVU) loop.
    * Manages high-level app states: `Home`, `Agent Creation`, `Active Chat`, and `Document Management`.
* **Component Library (`Bubbles`):**
    * Utilizes text inputs, progress bars for model downloads, and viewports for scrolling chat history.
* **Visual Design (`Lipgloss` & `Catwalk`):**
    * **Lipgloss:** Defines layouts, adaptive borders, and responsive padding.
    * **Catwalk:** Ensures consistent rendering across various terminal emulators and manages color palettes derived from specific LLM responses.
* **Rich Content Rendering (`Glow`):**
    * Full Markdown support for AI-generated reports, code blocks, and structured documentation.

---

### 2. The AI Gateway (Middleware)
This is the most critical layer for modularity. By using **Fantasy**, you are already positioned to handle structured data.

* **Agent Coordination:** Define agents as **Protocols**. An agent isn't just a prompt; it’s a set of capabilities (Tools + Model + State).
* **Provider Translation:** * Create a `Provider` interface in Go. This allows you to swap between `OpenAI`, `Anthropic`, or `Local` (Ollama/Llama.cpp) without changing the Gateway logic.
    * **Streaming Logic:** Ensure your Gateway supports **Server-Sent Events (SSE)** or equivalent streaming patterns so the TUI can update the chat bubbles in real-time as tokens arrive.

    This layer serves as the "brain," decoupling the UI from the underlying intelligence providers.

* **Agentic Framework (`Fantasy`):**
    * Manages the lifecycle of multi-agent coordination.
    * Enforces **Structured Responses** by mapping AI outputs to Go structs for predictable UI updates.
* **Provider Translation Layer:**
    * A unified interface to route requests to local (Ollama/LocalAI) or remote (OpenAI, Anthropic, Hugging Face) endpoints.
* **Model Management:**
    * Integration with **Hugging Face** for direct model downloading and local storage.
    * Support for user-provided GGUF or Safetensors files.
---

### 3. The Backend (Document & Multi-Modal Engine)
This layer handles the "heavy lifting" outside of the main UI loop.

#### **RAG & Document Management**
* **Vector Store:** For a local-first TUI, consider a lightweight vector store like **ChromaDB** or even a local **FAISS** index stored as a file.
* **PDF Processing:** Use a Go-native library like `unidoc/unipdf` to extract text. The Gateway then chunks this text and sends it to the embedding model.

#### **Audio & Communication (The "Omni" Features)**
To handle Whisper (STT), TTS, and Phone calls without locking the UI, these must run in **background Goroutines**.

| Feature | Implementation Strategy |
| :--- | :--- |
| **Whisper (STT)** | Capture audio via `portaudio` or `ffmpeg`, stream to a Whisper endpoint, and send a `WhisperCompleteMsg` back to Bubble Tea to populate the input field. |
| **TTS** | Stream response text to a TTS provider (like ElevenLabs or OpenAI) and play audio through a local buffer while the text renders in the TUI. |
| **Phone Calls** | Integrate **Twilio** or **SignalWire** APIs. The AI Gateway treats the phone call as a specialized "session" where input comes from a webhook rather than the keyboard. |
| **Tool Calling** | Use **JSON Schema** definitions. When the model requests a tool (e.g., `get_weather`), the Gateway executes the local Go function and feeds the result back to the LLM before the UI even sees the final answer. |

---

### Proposed State Machine Flow


1.  **Idle State:** User views the dashboard (New Chat / Load Session).
2.  **Configuration State (Huh):** User defines agent parameters.
3.  **Active Chat State (Bubble Tea):**
    * **Sub-state: Awaiting Input**
    * **Sub-state: Processing (Spinner component)**
    * **Sub-state: Streaming Response**
4.  **Document State (Glow):** User reviews a generated `DevSpec` or `CharterDoc`.
5.  **Action State:** A tool is triggered (e.g., making a phone call or searching a local PDF).

---

### Technical Tip: The "Msg" Pattern
In Bubble Tea, everything is a `Cmd`. To keep the backend clean, your Backend services should return functions that wrap the result in a custom message type:

```go
func (b *Backend) TranscribeAudio() tea.Cmd {
    return func() tea.Msg {
        text := b.WhisperClient.Run()
        return TranscriptionResultMsg(text)
    }
}
```
Ensure your TUI stays responsive (60fps) while the AI Gateway or Document Engine is doing complex processing in the background.

---

    

### 1. The Personality Schema
Instead of a generic system prompt, define a `Personality` struct. This allows the AI Gateway to "tune" the model’s behavior programmatically.

```go
type Personality struct {
    Name       string   `json:"name"`
    Archetype  string   `json:"archetype"` // e.g., "Stoic Mentor", "Sarcastic Tech Support"
    Traits     []string `json:"traits"`    // e.g., ["verbose", "witty", "cautious"]
    Speech     string   `json:"speech"`    // e.g., "Pirate slang", "Formal Victorian"
    CoreValues []string `json:"values"`    // Guides how the agent makes decisions
}
```

### 2. Updating the "Huh" Wizard
In your agent creation phase, add a specific **Persona Workshop**. Using `huh.Select` and `huh.MultiSelect`, you can let users build a personality without typing a wall of text.

* **Archetype Selection:** Use a list of presets (e.g., The Rebel, The Sage, The Assistant).
* **Tone Sliders:** While `huh` is largely text-based, you can use a `huh.Select` to set a "Temperament" (e.g., Aggressive <—> Passive).
* **Voice/Speech Toggles:** Use `huh.Confirm` to toggle specific quirks like "Uses Emojis" or "Never uses capital letters."

### 3. Implementing the "Persona Engine"
In your **AI Gateway (Middleware)**, you need a function that compiles these traits into a "Super-System Prompt." 

> **Tip:** Do not just append the traits. Use a template to give them weight.

```markdown
# Roleplay Protocol: {{.Name}}
You are currently roleplaying as {{.Name}}, a {{.Archetype}}.
## Constraints:
- Tone: {{.Traits}}
- Linguistic Style: {{.Speech}}
- Behavioral Bounds: {{.CoreValues}}

## Interaction Instructions:
Never break character. If the user asks about your AI nature, respond in the voice of {{.Name}}.
```

### 4. TUI Visual Cues (Lipgloss & Bubbles)
To make roleplaying feel immersive in a terminal, the interface should visually reflect who you are talking to.

* **Dynamic Borders:** Use `Lipgloss` to change the border color of the chat bubble based on the agent's personality. 
    * *Example:* A "Hot-Headed" agent gets a `Red` border; a "Calm" agent gets `ForestGreen`.
* **Avatar Rendering:** Since you are in the terminal, use **ASCII Art** or **ANSI icons** in the chat header. You can store a small library of ASCII avatars that match the `Archetype` selected in the Wizard.
* **Status Bars:** Use a bottom status bar to show the agent's current "Mood" or "State" if your state machine tracks emotional shifts during the conversation.

### 5. Advanced: Memory & Continuity
For deep roleplay, the agent needs to remember past interactions.
* **The "Vibe" Buffer:** In your **Backend (Document Engine)**, store a "Memory Summary" in your RAG system. Instead of just searching for facts, search for "Relational Data"—how the agent feels about the user based on previous chats.
* **Fantasy Integration:** Use Fantasy's `Structured Response` to have the model output two things: the **In-Character Response** and a hidden **Emotional State Update** that you save to the database.

---

### Example Archetype: The Roman Senator
Since you appreciate historical accuracy, you could have a preset for a **Roman Senator (90 BCE)**:
* **Traits:** ["Authoritative", "Traditionalist", "Oratorical"]
* **Speech:** "Latinate sentence structures, refers to the Mos Maiorum frequently."
* **TUI Style:** Purple (`latus clavus`) borders and a marble-white text color.

---

Build a Chat TUI using the Charmbracelet family of repositories. Use  the Fantasy framework. https://github.com/charmbracelet/fantasy for the AI layer. Use BubbleTeam for the TUI Framework  https://github.com/charmbracelet/bubbletea for the TUI Framework. Use Lipgloss https://github.com/charmbracelet/lipgloss for the layout. Use Glow https://github.com/charmbracelet/glow for the markdown renderer for documents. Use Catwalk https://github.com/charmbracelet/catwalk for a database of CHarmbracelet compatible LLM models.  Use Bubbles https://github.com/charmbracelet/bubbles for more TUI components. And use Huh https://github.com/charmbracelet/huh for the selector tools in the beginning and for the wizard. 


Technical Architecture
To ensure scalability and clean code, we will divide the application into three primary layers:
1. The Frontend (TUI Layer)
Huh & Bubbles: Used for the agent creation Wizard" phase. Huh handles walking through user creating an agent and the various settings
Bubble Tea: The main orchestration engine. It will manage the main chatUI interface (New Chat, New Agent, Upload PDF)
Lipgloss: Used to define the visual style, borders, and padding for a professional TUI look.
Glow & Catwalk: Glow will render the generated Markdown docs  Catwalk ensures the TUI remains compatible with various terminal themes and LLM-derived color palettes.
2. The AI Gateway (Middleware)
Fantasy Framework: Acts as the "brain." It will handle providers, endpoint mapping, agent coordination, etc. Structured Response capability will define models for the AI to follow. 
Provider Translation: A distinct service that takes the internal request and routes it to the correct endpoint (Local, OpenAI, Anthropic) using keys dedicated settings screen.
3. The Backend (Document Engine)
State Machine: Tracks which state the user is in
RAG and Question Ansering
support document and PDF upload, chat with PDF, and Question Answering
Support transcrobe via whispr
Support TTS
Support outbound phone call making
Support tool calling
support mcp servers and websocket or a2a protocol to faciliate agentic communication. the chatroom should be stateful so that the agent doesnt need to pass the full conversation with each pass, it is stored in the chatroom. support prompt processing, reasoning, tool use and function calling. support downloading models from hugging face in system or allow user to upload  model. support text embedding

---
---
### 1. The "Tooey" Dashboard (A Multi-Pane Workspace)
**Tooey** is CharmingMan's most distinct "think-big" feature. Instead of just chatting, it acts as a **parallel reasoning engine**:
* **Split-Screen View:** You can open multiple chat windows side-by-side. This allows you to compare different models (e.g., GPT-5 vs. Claude 4) or different personas in real-time.
* **Drawer System:** A left-hand drawer houses your conversation history and folders, while a right-hand "Tools" or "Settings" panel can be toggled for quick model adjustments.

* **Resizable Panes:** Like a code editor, panes can be resized or collapsed to focus on a specific workflow.

### 2. Persona-Driven Architecture
The dashboard isn't just a text box; it's a **Persona Hub**.
* **Advanced Personas:** You can switch between "Expert" personas that come pre-loaded with specific system instructions, model preferences, and even attached files.
* Roadmap
* **Visual Identity:** Each persona can have a distinct visual style and color palette, helping you keep track of which "expert" you are currently consulting.

### 3. Technical Dashboard Features
* **Live File Sync:** A dedicated area for managing local files. You can drop PDFs, code, or images directly into the workspace, and the "Live Sync" feature can track changes on your hard drive in real-time.
* **Markdown & Diagrams:** The main workspace features an integrated **Glow** renderer. It doesn't just show text; it renders complex Mermaid diagrams, PlantUML, and code blocks with syntax highlighting as the model "thinks."
* Roadmap:
* **Cost & Token Monitoring:** A transparent reporting area shows real-time token usage and costs across different API providers, giving you a "cloud-billing" style overview of your AI spend.

### 4. Interaction Modes
* **Chat Mode:** A traditional terminal chat UI with all the main Tooey features.
* Roadmap:
* **Call Mode:** A dedicated voice-first interface within the dashboard for hands-free interaction.
* **Draw Mode:** An integrated image generation workspace (utilizing DALL-E 3 or local diffusion models) that lives alongside your text chats.
* **Search Mode:** A dynamic web-filtering toggle that allows the workspace to pull in real-time data via Google Search or YouTube transcripts.
* **Scatter/Gather Workflow:** You send one prompt to multiple models at once (Scatter). They generate answers in parallel rays.
* **Consensus & Fusion:** The dashboard displays these responses side-by-side. You can then use a "Fusion" algorithm to merge the best parts of each answer into one "gold" response (Gather), effectively using the models to fact-check each other and eliminate hallucinations.

## The CharmingMan "Command Center" Design

### 1. The Multi-Pane "Tiled" Workspace
This main dashboard uses a **Flexbox-style layout** powered by `Lipgloss`.

* **The Quad-View:** Support for up to four active chat panes simultaneously. Each pane is an independent `Bubble Tea` model.
* **Active Focus:** The "focused" pane is highlighted with a bold, `thicker` border and a specific color (e.g., Mauve for GPT-4, Green for local Ollama) using `Catwalk` to ensure the colors pop against the terminal background.
* **The "Beam" TUI Implementation:** A dedicated horizontal split where the top 70% shows the multi-model responses side-by-side, and the bottom 30% is a shared input buffer for "scattering" the prompt.


### 2. The Sidebar "Intelligence Hub"
A collapsible left-hand navigation bar (toggled with `ctrl+b`) manages the application's global state.

* **Agent Library:** A vertical list of custom agents (created via the `Huh?` wizard). Icons next to their names indicate the backend (e.g., a "☁️" for Anthropic, a "🏠" for local Llama).
* **Active MCP Servers:** A status indicator showing which Model Context Protocol servers are currently attached (e.g., "Filesystem: Active", "PostgreSQL: Connected").
* **Token Telemetry:** A sparkline-style mini-graph in the bottom corner of the sidebar showing real-time token throughput or cost-per-session.

---

### 3. The "Stage" & Document Engine
The right-hand side of the dashboard is reserved for **The Stage**, a context-aware viewport for non-chat artifacts.

* **Document Preview:** When you "Chat with a PDF," the document text is rendered in this pane using `Glow`. As the AI references specific pages, the `Glow` viewport auto-scrolls to the relevant section.
* **Code & Artifacts:** If the agent generates code or a "DevSpec," it is promoted to The Stage. This keeps the chat history clean while allowing you to inspect the output in full-screen width.
* **Live Reasoning Log:** A "Thinking" drawer that shows the hidden chain-of-thought, tool-calling logs, and search results without cluttering the primary conversation.

---

### 4. Interactions & Navigation
The dashboard prioritizes keyboard-driven speed, mimicking the efficiency of an IDE like Neovim.

| Keybinding | Action | Component |
| :--- | :--- | :--- |
| `Tab` | Cycle focus between active chat panes | `Bubble Tea` |
| `Ctrl + N` | Open the `Huh?` Wizard for a new Agent/Chat | `Huh?` |
| `Ctrl + S` | Toggle the "Stage" (Document/Glow view) | `Lipgloss` |
| `/` | Enter "Command Mode" for tool calling or model switching | `Bubbles` |

---

### 5. Visual Aesthetic (The "Big-AGI" Polish)
To achieve that professional, high-end feel in a terminal:
* **Adaptive Theme:** Use `Catwalk` to sync the TUI theme with the user's terminal (Dracula, Solarized, Nord) while injecting specific accent colors for different AI providers.
* **Status Lines:** A persistent bottom bar (similar to `vim-airline`) showing the current model, temperature, active "tools" (e.g., Whisper, MCP), and system latency.
* **Glow-Enhanced Response:** AI responses aren't just strings; they are processed through `Glow` to provide syntax highlighting for code and beautiful tables for structured data.

Other Screens:
---

### 1. The "Canvas" (Node-Graph Layout)
Instead of linear panes, the workspace acts like a **tiling window manager** (like `iTerm2` splits or `tmux`), but with a focus on "pipelines."
* **The Visual:** A tree-like structure where each box is a `Bubble Tea` component. 
* **The Flow:** User Input → Coordinator Agent → (Splits into 3 Specialist Agents) → Aggregator Agent → Final Glow Markdown.
* **Interaction:** You use the arrow keys to "navigate" the graph. Highlighting a node expands it to show that agent's specific logs, while the rest of the nodes shrink to "sparkline" summaries using `Lipgloss` borders.

### 2. The "Broadcast" (Central Stage + Ghost Peers)
Inspired by **big-AGI’s Beam**, but optimized for a "Manager" view.
* **Main Stage:** A large central area (80% width) for the "Lead Agent" or the final consolidated response rendered in `Glow`.
* **The Ghost Stream:** A narrow right-hand column showing real-time, high-speed "ghost" text from sub-agents working in the background. It looks like a scrolling log (tail -f style). 
* **Utility:** If a sub-agent hits an error or needs "Human-in-the-loop" (HITL) approval, that "ghost" pane flashes red and expands, demanding focus.

### 3. The "Focus Mode" (Zen/Modal Layout)
For deep reasoning or long-form document generation (like your **DevSpecs**).
* **Minimalist UI:** All sidebars and status lines auto-hide (`ctrl+z`). The terminal becomes a single, centered column of text with generous `Lipgloss` padding.
* **Breadcrumb Navigation:** A single, thin line at the top showing the path: `Project > Sprint 1 > Architecture_Agent > DevSpec_v2`.
* **HUD Overlay:** A transparent-looking head-up display (using high-contrast background colors) that only appears when you trigger a command (`/`) or search for an MCP tool.

### 4. The "Dashboard" (Executive Overview)
Perfect for when you have multiple persistent agentic "rooms" running in the background.
* **Grid System:** A 2x2 or 3x3 grid of mini-viewports.
* **Content:** Each viewport shows a "State Card":
    * **Agent Name:** (e.g., "History_Expert")
    * **Status:** (Idle, Thinking, or Calling Tool)
    * **Context usage:** (e.g., 85% of 128k tokens)
    * **Last Action:** "Searched: Battle of Pydna"
* **Interaction:** Pressing a number (1-9) "zooms" into that specific agent's full chat room.

### 5. The "Notebook" (Linear Iterative)
Similar to Jupyter or Claude's "Artifacts" side-by-side.
* **Left Side:** A persistent, vertical "history" of code blocks, prompts, and tool calls.
* **Right Side:** The "Outcome" area. Every time the code/prompt on the left is updated and run, the right side refreshes. 
* **Usage:** Ideal for the "Spec Driven Design" tool you're building, where you can tweak a requirement on the left and see the generated PromptDoc update in real-time on the right.

---

### Suggested Layout Mapping for CharmingMan

| Task | Recommended Layout | Charm Component |
| :--- | :--- | :--- |
| **Agent Creation** | **Wizard/Focus** | `Huh?` |
| **Multi-Model Comparison** | **Tiled (Beam)** | `Bubble Tea` + `Catwalk` |
| **Complex Multi-Agent Task** | **Canvas/Node** | `Lipgloss` (Flexbox) |
| **Long-form Doc Review** | **Notebook** | `Glow` + `Bubbles` (Viewport) |


* **The Markdown-Stage:** A `Glow` powered area for the final output.
* **The Tool-Belt:** A grid of buttons (clickable via mouse) to trigger specific MCP tools or local Python scripts via `uv`.

### 3. Sharing via "Charming Registry"
Imagine a command like `charmingman ui install big-agi-clone`. 
* **Implementation:** The app fetches the manifest from a GitHub repo or a central registry.
* **Theming:** Using **Catwalk**, the shared UI automatically adapts to the user's local terminal colors, so a UI designed in "Dracula" still looks great for a user on "Solarized Light."

### 4. The "Live Designer" Mode
Since you mentioned mouse support (resizing/moving), you could have a **"Design Mode"** within the app:
Yes, **Bubble Tea** (the core engine of the Charmbracelet family) has comprehensive mouse support, but there is a catch: it provides the **building blocks**, not a pre-built "Window Manager."

While there isn't a single "Window" component you can drop in that just works like macOS or Windows, you can absolutely build one using their event system.

### How it works technically:
1.  **Enabling Mouse Support:** You have to explicitly tell your program to listen for mouse events by setting the `MouseMode` in your `View()` (e.g., `tea.MouseModeAllMotion`).
2.  **Tracking State:** You define a "window" as a struct with an `X`, `Y`, `Width`, and `Height`.
3.  **Handling Events:**
    * **Moving:** When a `tea.MouseClickMsg` hits your "title bar" area, you set a `dragging` flag to true. As long as the button is held, you use `tea.MouseMotionMsg` to update the window's `X` and `Y` coordinates based on the cursor's movement.
    * **Resizing:** Similar to moving, but you listen for clicks on the border or a corner "handle." When the mouse moves, you increment/decrement the `Width` and `Height` values of your component.
4.  **Rendering:** You use **Lipgloss** to render the window at those specific coordinates using `lipgloss.Place()` or by calculating the padding/margins dynamically.

### Is there a shortcut?
If you want to avoid writing the math for moving and resizing from scratch, keep an eye on these specific Charm projects:
* **`Lipgloss v2`:** Recently introduced a `Canvas` and `Layer` system that makes "Z-index" (stacking windows on top of each other) and coordinate-based positioning much easier than the original version.
* **The "Bubbles" Viewport:** While primarily for scrolling, it’s often used as the "inner pane" of a custom window.

### For CharmingMan:
Since you’re going for a **big-AGI** style dashboard, using the mouse to resize panes is a killer feature for a TUI. It makes the "terminal" feel like a high-end IDE. You could even implement "double-click to maximize" a pane, which is a common power-user move in those layouts.

**Pro-tip:** Most TUI users still prefer keyboard shortcuts, so if you build mouse-resizable windows, definitely map `Ctrl+Arrows` or similar to do the same thing for the "no-mouse" crowd!



