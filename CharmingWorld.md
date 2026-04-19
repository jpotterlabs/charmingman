# CharmingWorld: Roleplay & Simulation Engine

**CharmingWorld** is the simulation layer of CharmingMan. It leverages the underlying multi-agent orchestration to create immersive, roleplay-driven environments.

## 1. Personality Schema

Instead of generic system prompts, CharmingWorld uses a structured `Personality` definition. This allows the AI Gateway to tune model behavior programmatically.

- **Archetypes**: Presets like "Stoic Mentor," "Sarcastic Tech Support," or "Roman Senator."
- **Speech Style**: Define linguistic rules (e.g., "Uses pirate slang" or "Formal Victorian").
- **Core Values**: Guides how the agent makes decisions and stays in character.

## 2. World Book (The Codex)

The World Book is a secondary RAG (Retrieval-Augmented Generation) layer specifically for setting-specific facts, history, and rules.

### Features:
- **Keyword-Triggered Injection**: The backend scans user messages for keywords (e.g., "The Forum") and dynamically injects the relevant lore snippet into the prompt.
- **Lore Entry Management**: Create and edit entries using a specialized `Huh?` form.
- **Lore Discovery**: Some entries are hidden and only become "discovered" (visible in the TUI) once the agent mentions them in conversation.

## 3. Immersive Interface

The **Infinity Canvas** and **YAML Layouts** play a crucial role in creating an immersive experience:

- **Dynamic Visual Cues**: Window borders and accent colors change based on the agent's personality and "mood."
- **Spatial Lore Graphs**: Use the Infinity Canvas to map out connections between characters and locations.
- **Character Avatars**: High-quality ASCII or ANSI icons in window headers to visually distinguish between different agents.

## 4. Continuity & Memory

- **The "Vibe" Buffer**: Store a memory summary in the RAG system to track "Relational Data"—how the agent feels about the user based on previous interactions.
- **Stateful Chatrooms**: Persistent **RoomID** context ensures that multiple agents interacting in the same room maintain a consistent shared history.

---
*CharmingWorld elevates CharmingMan from a chat tool to a living, breathing terminal simulation.*
