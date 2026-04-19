# CharmingMan Usage Guide 🎩

This guide provides a walkthrough for using the **CharmingMan** multi-agent swarm, including text and voice-based interactions.

## 🎙️ Voice Input (The 'v' Key)

The most advanced feature of CharmingMan is hands-free interaction via **Whisper STT**.

### 1. Prerequisites for Voice
To use voice input, ensure you have the following:
- **Audio Recorder**: `sox` must be installed on your system.
    - macOS: `brew install sox`
    - Ubuntu/Debian: `sudo apt-get install sox`
- **Whisper API Key**: The AI Gateway must be configured with a valid `OPENAI_API_KEY` in its `.env` file.

### 2. Using Voice in the TUI
- **Activate Recording**: While in the Chat or Infinity Canvas view, press the **'v'** key.
- **Recording Duration**: The system will record for **3 seconds** (current default for testing). A spinner will appear to indicate recording is in progress.
- **Auto-Transcription**: Once finished, the audio is sent to the AI Gateway's `/api/v1/transcribe` endpoint.
- **Automatic Routing**: The transcribed text is automatically routed to the primary agent as a new prompt.

---

## 💬 Multi-Agent Swarm Interaction

CharmingMan allows you to coordinate multiple AI agents simultaneously within a shared context.

### 1. @Mention Routing
You can direct messages to specific agents using the `@Name` syntax:
- `@Architect How should I structure this Go project?`
- `@Coder Implement a Bubble Tea model for the home screen.`

If no mention is used, the message is sent to the **primary agent** (usually the first one created).

### 2. Shared Room Context
All agents in a session share a **RoomID**. This means:
- Agents are aware of previous messages from both you and other agents.
- You can say: `@Coder Look at what the @Architect said and implement the first step.`

### 3. Navigation Controls
- **TAB**: Cycle focus between active agent panes.
- **Mouse Drag**: Move agent windows anywhere on the **Infinity Canvas**.
- **Ctrl+Scroll / +/-**: Zoom the camera in and out of the canvas.
- **Arrow Keys / WASD**: Pan the camera across the workspace.

---

## 📚 Grounding with RAG

To make your agents smarter, provide them with your own data:
1. Place your documents (`.pdf`, `.md`, `.txt`) in the `documents/` directory.
2. The AI Gateway will automatically chunk and index these files.
3. In your prompt, use the `-rag` flag (or enable RAG in the agent settings) to trigger context injection:
   - `Explain the project structure -rag`

---

## 🛠️ Configuration Recap

| Feature | Key Requirement |
|---------|-----------------|
| **Voice Input** | `sox` installed + `OPENAI_API_KEY` |
| **RAG** | `OPENAI_API_KEY` (for embeddings) |
| **Claude Models** | `ANTHROPIC_API_KEY` |
| **Local Models** | Ollama running on `localhost:11434` |
