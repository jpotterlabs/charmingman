# Multimodal: Voice & Speech

CharmingMan is not limited to text. It includes a multimodal suite that allows for voice-based interactions and speech synthesis for agent responses.

## 🎙️ Voice Input (Whisper STT)

You can prompt your agents using your voice. The TUI captures audio from your microphone, sends it to the AI Gateway, and automatically routes the transcribed text to your agents.

### How it works:
1. **Trigger**: Press the `'v'` key in the TUI to start recording.
2. **Capture**: The TUI calls the system's `rec` command (from the `sox` package) to capture 3 seconds of audio.
3. **Transcribe**: The audio file is sent to the AI Gateway's `POST /api/v1/transcribe` endpoint.
4. **Processing**: The Gateway uses **OpenAI Whisper** to convert the audio to text.
5. **Routing**: The resulting text is sent back to the TUI and automatically routed to the primary chat agent as a new prompt.

### Prerequisites
To use voice input, you must have `sox` installed on your system:
- **macOS**: `brew install sox`
- **Linux**: `sudo apt-get install sox`

## 🔊 Agent Speech (TTS)

You can listen to agent responses instead of just reading them.

### How it works:
1. **Trigger**: Press the `'s'` key when focused on an agent's response.
2. **Generate**: The TUI calls the AI Gateway's `POST /api/v1/speech` endpoint with the response text.
3. **Stream**: The Gateway uses **OpenAI TTS** to generate an MP3 audio stream.
4. **Playback**: The TUI receives the stream, saves it to a temporary file, and plays it back using the system's `play` command.

## ⚙️ Configuration

Set up multimodal features in your `.env`:
```env
OPENAI_API_KEY=sk-... # Required for Whisper and TTS
```

## 🚀 Key Features

- **Hands-free Prompts**: Ideal for quick queries or brainstorming while away from the keyboard.
- **Natural Interaction**: Makes the agents feel more alive and accessible.
- **Cross-Platform**: Leverages standard CLI tools (`sox`) for broad compatibility.

## 🛠️ Implementation Details

- **TUI**: `internal/tui/audio.go` manages the recording state and spinner animation.
- **Backend**: `backend/internal/handler/transcribe.go` handles the interface with OpenAI's audio services.
