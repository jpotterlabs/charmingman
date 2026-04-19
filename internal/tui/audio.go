package tui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type transcribeMsg struct {
	Text  string
	Error error
}

type recordingFinishedMsg struct {
	Path  string
	Error error
}

type VoiceInputModel struct {
	spinner   spinner.Model
	recording bool
	audioPath string
	err       error
	done      bool
	width     int
	height    int
}

func NewVoiceInputModel() *VoiceInputModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	
	return &VoiceInputModel{
		spinner: s,
	}
}

func (m *VoiceInputModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m *VoiceInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyPressMsg:
		switch msg.String() {
		case "v": // Start recording on 'v'
			if !m.recording {
				m.recording = true
				m.err = nil
				m.done = false
				return m, m.recordAudio()
			}
		case "esc", "q":
			m.done = true
		}

	case recordingFinishedMsg:
		m.recording = false
		if msg.Error != nil {
			m.err = msg.Error
			return m, nil
		}
		m.audioPath = msg.Path
		return m, m.transcribeAudio(msg.Path)

	case transcribeMsg:
		m.done = true
		if msg.Error != nil {
			m.err = msg.Error
			return m, nil
		}
		// In a real app, we would send this text to the chat input
		return m, func() tea.Msg { return RouteMsg{Mention: "primary-chat", Prompt: msg.Text} }

	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *VoiceInputModel) View() tea.View {
	if m.recording {
		return tea.NewView(fmt.Sprintf("%s Recording... Press Ctrl+C to stop (simulated)", m.spinner.View()))
	}
	if m.err != nil {
		return tea.NewView(fmt.Sprintf("Error: %v", m.err))
	}
	if m.done {
		return tea.NewView("Transcription complete.")
	}
	return tea.NewView("Press 'v' to start voice input.")
}

func (m *VoiceInputModel) recordAudio() tea.Cmd {
	return func() tea.Msg {
		// Simplified recording logic: use 'sox' or 'ffmpeg' to record 3 seconds
		// In a production app, we'd use a proper portaudio/pulse bind
		path := filepath.Join(os.TempDir(), "charmingman_record.wav")
		
		// Attempt to record for 3 seconds using sox (rec)
		cmd := exec.Command("rec", "-r", "16000", "-c", "1", path, "trim", "0", "3")
		err := cmd.Run()
		
		return recordingFinishedMsg{Path: path, Error: err}
	}
}

func (m *VoiceInputModel) transcribeAudio(path string) tea.Cmd {
	return func() tea.Msg {
		gatewayURL := os.Getenv("GATEWAY_URL")
		if gatewayURL == "" {
			gatewayURL = "http://localhost:8090/api/v1/transcribe"
		} else {
			u, _ := url.Parse(gatewayURL)
			u.Path = "/api/v1/transcribe"
			gatewayURL = u.String()
		}

		file, err := os.Open(path)
		if err != nil {
			return transcribeMsg{Error: err}
		}
		defer file.Close()

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", filepath.Base(path))
		if err != nil {
			return transcribeMsg{Error: err}
		}
		io.Copy(part, file)
		writer.Close()

		req, err := http.NewRequest("POST", gatewayURL, body)
		if err != nil {
			return transcribeMsg{Error: err}
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("X-Charming-Key", os.Getenv("GATEWAY_API_KEY"))

		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return transcribeMsg{Error: err}
		}
		defer resp.Body.Close()

		var result struct {
			Text string `json:"text"`
		}
		json.NewDecoder(resp.Body).Decode(&result)

		return transcribeMsg{Text: result.Text}
	}
}
