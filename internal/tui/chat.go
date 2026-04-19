package tui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type chatResponseMsg struct {
	Response string
	Usage    interface{}
	Error    error
}

type ChatModel struct {
	viewport  viewport.Model
	textinput textinput.Model
	history   []string
	ready     bool
	width     int
	height    int
	Focused   bool

	// Temporary storage for agent config
	Provider string
	Model    string
	APIKey   string
}

func NewChatModel() *ChatModel {
	ti := textinput.New()
	ti.Placeholder = "Send a message..."

	return &ChatModel{
		history:   make([]string, 0),
		textinput: ti,
	}
}

func (m *ChatModel) Init() tea.Cmd {
	return nil
}

func (m *ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Subtract space for text input
		viewportHeight := m.height - 3
		if viewportHeight < 0 {
			viewportHeight = 0
		}

		if !m.ready {
			m.viewport = viewport.New(viewport.WithWidth(m.width), viewport.WithHeight(viewportHeight))
			m.ready = true
		} else {
			m.viewport.SetWidth(m.width)
			m.viewport.SetHeight(viewportHeight)
		}

	case tea.KeyPressMsg:
		if m.Focused {
			switch msg.String() {
			case "enter":
				userInput := m.textinput.Value()
				if userInput != "" {
					m.AddMessage(fmt.Sprintf("You: %s", userInput))
					m.textinput.SetValue("")
					cmds = append(cmds, m.sendMessage(userInput))
				}
			}
			m.textinput, cmd = m.textinput.Update(msg)
			cmds = append(cmds, cmd)
		}

	case chatResponseMsg:
		if msg.Error != nil {
			m.AddMessage(fmt.Sprintf("Error: %v", msg.Error))
		} else {
			m.AddMessage(fmt.Sprintf("AI: %s", msg.Response))
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *ChatModel) View() tea.View {
	if !m.ready {
		return tea.NewView("Initializing...")
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		m.viewport.View(),
		"\n",
		m.textinput.View(),
	)
	return tea.NewView(content)
}

func (m *ChatModel) AddMessage(msg string) {
	m.history = append(m.history, msg)
	m.viewport.SetContent(strings.Join(m.history, "\n"))
	m.viewport.GotoBottom()
}

func (m *ChatModel) SetFocused(f bool) {
	m.Focused = f
	if f {
		m.textinput.Focus()
	} else {
		m.textinput.Blur()
	}
}

func (m *ChatModel) sendMessage(prompt string) tea.Cmd {
	return func() tea.Msg {
		url := os.Getenv("GATEWAY_URL")
		if url == "" {
			url = "http://localhost:8090/api/v1/chat"
		}
		
		payload := map[string]string{
			"provider": m.Provider,
			"model":    m.Model,
			"prompt":   prompt,
		}
		
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return chatResponseMsg{Error: err}
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			return chatResponseMsg{Error: err}
		}

		req.Header.Set("Content-Type", "application/json")
		// Use the API key if provided, or fallback to a default for local testing
		apiKey := m.APIKey
		if apiKey == "" {
			apiKey = "charming-secret-key"
		}
		req.Header.Set("X-Charming-Key", apiKey)

		client := &http.Client{Timeout: 35 * time.Second} // Slightly longer than gateway timeout
		resp, err := client.Do(req)
		if err != nil {
			return chatResponseMsg{Error: err}
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			var errResp struct {
				Error string `json:"error"`
			}
			json.NewDecoder(resp.Body).Decode(&errResp)
			return chatResponseMsg{Error: fmt.Errorf("server error (%d): %s", resp.StatusCode, errResp.Error)}
		}

		var result struct {
			Response string      `json:"response"`
			Usage    interface{} `json:"usage"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return chatResponseMsg{Error: err}
		}

		return chatResponseMsg{
			Response: result.Response,
			Usage:    result.Usage,
		}
	}
}
