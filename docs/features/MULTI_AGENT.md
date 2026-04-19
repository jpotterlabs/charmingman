# Multi-Agent Orchestration

CharmingMan is designed from the ground up as a **Multi-Agent Simulation Environment**. Unlike single-agent assistants, CharmingMan allows you to create, manage, and coordinate a "swarm" of agents, each with distinct personas and specialized models.

## 🧠 Core Concepts

### 1. Agent Personas
Each agent is defined by a manifest that includes:
- **Name**: A unique identifier for @mentions.
- **Model & Provider**: The underlying LLM (e.g., GPT-4o on OpenAI or Llama3 on Ollama).
- **Persona**: A system prompt that defines the agent's behavior, expertise, and tone.
- **Tools**: A list of capabilities (via MCP) the agent is authorized to use.

### 2. @Mention Routing
Communication in CharmingMan is handled through a broadcast system. When you type a message:
- **Standard Send**: Routes the prompt to the currently focused chat window.
- **@Name Mention**: You can route a prompt to any agent on the canvas by prefixing your message with `@AgentName`.

**Technical Flow:**
1. The TUI `ChatModel` detects an `@` prefix in the input.
2. It wraps the prompt into a `RouteMsg{Mention: "AgentName", Prompt: "..."}`.
3. The `Manager` in `internal/tui/manager.go` broadcasts this message to all windows.
4. The window whose ID matches the `Mention` field processes the prompt and calls the AI Gateway.

### 3. Shared Room Context
One of CharmingMan's most powerful features is the **Shared Room Context**. All agents active in a workspace share a common `RoomID` (a UUID generated when the session starts).

- **Unified History**: Every message (from you or any agent) is persisted in the AI Gateway's database under this `RoomID`.
- **Collaborative Reasoning**: When an agent is called, the Gateway retrieves the last 10 messages from that room. This means agents can see the responses of *other* agents, enabling multi-step workflows and collaborative problem-solving.

```mermaid
sequenceDiagram
    participant User
    转向 TUI[Manager]
    participant A1[Agent: Alice]
    participant A2[Agent: Bob]
    participant GW[AI Gateway]
    
    User->>TUI: "@Alice help me plan the project"
    TUI->>A1: RouteMsg(Mention: Alice)
    A1->>GW: POST /chat (Prompt, RoomID: X)
    GW-->>A1: "I've outlined the tasks."
    A1->>TUI: Display response
    
    User->>TUI: "@Bob review Alice's plan"
    TUI->>A2: RouteMsg(Mention: Bob)
    A2->>GW: POST /chat (Prompt, RoomID: X)
    Note over GW: GW fetches history for Room X (includes Alice's response)
    GW-->>A2: "The plan looks solid, but needs more testing."
    A2->>TUI: Display response
```

## 🛠️ Usage in TUI

- **Creation**: Use the interactive Wizard on startup to configure your initial agents.
- **Switching Focus**: Press `Tab` to cycle between agent windows.
- **Coordination**: Use `@AgentName` in the input field of *any* agent to trigger a different one.

## 🚀 Key Advantages

- **Model Heterogeneity**: Run a heavy GPT-4o model for planning and a fast, local Llama3 model for summaries in the same room.
- **Contextual Persistence**: No need to copy-paste between windows; the room history handles the handoff.
- **Human-in-the-Loop**: You act as the conductor, deciding which agent speaks next and when to pivot the conversation.
