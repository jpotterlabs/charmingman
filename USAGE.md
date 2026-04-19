# User Guide: Interacting with your Swarm

CharmingMan is a multi-agent environment where you can coordinate several AI assistants at once. This guide covers how to use the TUI features effectively.

## 🌌 Canvas Navigation

Explore your spatial workspace using these shortcuts:

| Command | Action |
|---------|--------|
| **Arrow Keys** | Pan the camera (move the viewport) |
| **Drag Background** | Pan the camera (using left mouse) |
| **`+` / `-`** | Zoom in or out |
| **`0`** | Reset camera and zoom |
| **`Tab`** | Cycle focus between open windows |
| **Drag Title Bar**| Move a window around the canvas |
| **Drag Corner** | Resize a window (bottom-right handle) |

## 💬 Multi-Agent Interaction

### Direct Prompts
To talk to the focused agent, just type in the input field and press `Enter`.

### @Mentions
You can trigger *any* agent from *any* input field by using an `@mention`.
- **Syntax**: `@AgentName your message`
- **Example**: `@Alice what do you think of this?`

The message will be routed specifically to Alice, even if you are currently focused on a different agent's window.

### Shared History
All agents in your room see each other's responses. You can have one agent plan a task and another agent review that plan automatically.

## 🎙️ Voice & Speech

### Voice Input (STT)
1. Ensure `sox` is installed on your system.
2. Press `'v'` in the TUI to start recording.
3. Speak your prompt. Transcription will begin automatically after 3 seconds.
4. The transcribed text will be routed to your primary agent.

### Agent Speech (TTS)
1. Ensure `sox` is installed on your system.
2. Focus on an agent's response in the chat history.
3. Press `'s'` to generate and play audio for that response.

## 🛠️ MCP Tools

View your available agent capabilities in the **Tool Belt** window.
- Agents can automatically decide to use these tools (e.g., searching a filesystem) based on your prompt.
- Press `'r'` in the Tool Belt window to refresh the tool list from the AI Gateway.

## 📚 Grounding with RAG

Agents can use documents you've uploaded to provide more accurate answers.
- The chat history will indicate when sources were found: `(Sources: 3 found)`.
- Use **"The Stage"** window to preview document content and RAG snippets.
