package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"charmingman/internal/tui"
	tea "charm.land/bubbletea/v2"
)

type state int

const (
	stateWizard state = iota
	stateDashboard
	stateSavingAgent
)

type agentSavedMsg struct {
	Error error
}

type rootModel struct {
	state   state
	wizard  *tui.WizardModel
	manager *tui.Manager
	width   int
	height  int
}

func initialModel() rootModel {
	return rootModel{
		state:   stateWizard,
		wizard:  tui.NewWizardModel(),
		manager: tui.NewManager(),
	}
}

func (m rootModel) Init() tea.Cmd {
	return m.wizard.Init()
}

func (m rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case stateWizard:
		newWizard, cmd := m.wizard.Update(msg)
		m.wizard = newWizard.(*tui.WizardModel)
		if m.wizard.IsDone() {
			m.state = stateSavingAgent
			return m, m.saveAgent()
		}
		return m, cmd

	case stateSavingAgent:
		if msg, ok := msg.(agentSavedMsg); ok {
			if msg.Error != nil {
				// For now, just log and continue to dashboard anyway
				// In a real app, we might show an error screen
				fmt.Printf("Error saving agent: %v\n", msg.Error)
			}
			m.state = stateDashboard
			return m, m.initDashboard()
		}
		return m, nil

	case stateDashboard:
		switch msg := msg.(type) {
		case tea.KeyPressMsg:
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "tab":
				if len(m.manager.Windows) > 1 {
					focusedIdx := -1
					for i, w := range m.manager.Windows {
						if w.Focused {
							focusedIdx = i
							break
						}
					}
					// Focus next window
					nextIdx := (focusedIdx + 1) % len(m.manager.Windows)
					m.manager.FocusWindow(m.manager.Windows[nextIdx].ID)
				}
			}
		case tea.MouseMsg:
			cmd := m.manager.HandleMouse(msg)
			if cmd != nil {
				return m, cmd
			}
		case tea.WindowSizeMsg:
			m.width = msg.Width
			m.height = msg.Height
			m.manager.SetSize(msg.Width, msg.Height)
		}
		cmd := m.manager.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m rootModel) View() tea.View {
	var v tea.View
	v.AltScreen = true
	v.MouseMode = tea.MouseModeAllMotion

	switch m.state {
	case stateWizard:
		v.SetContent(m.wizard.View().Content)
	case stateSavingAgent:
		v.SetContent("Saving agent to gateway...")
	case stateDashboard:
		v.SetContent(m.manager.View().Content)
	}
	return v
}

func (m rootModel) saveAgent() tea.Cmd {
	return func() tea.Msg {
		url := os.Getenv("GATEWAY_URL")
		if url == "" {
			url = "http://localhost:8090/api/v1/agents"
		} else {
			// Ensure it points to the agents endpoint if it was just the base URL
			if !strings.HasSuffix(url, "/agents") {
				url = strings.TrimSuffix(url, "/chat") + "/agents"
			}
		}

		// Map model names to providers for the backend
		provider := "openai"
		model := m.wizard.Results.Model
		switch {
		case strings.Contains(strings.ToLower(model), "gpt"):
			provider = "openai"
		case strings.Contains(strings.ToLower(model), "claude"):
			provider = "anthropic"
		case strings.Contains(strings.ToLower(model), "llama3"):
			provider = "ollama"
		}

		payload := map[string]string{
			"name":     m.wizard.Results.Name,
			"model":    model,
			"provider": provider,
			"persona":  m.wizard.Results.Persona,
			"api_key":  m.wizard.Results.APIKey,
		}

		jsonData, err := json.Marshal(payload)
		if err != nil {
			return agentSavedMsg{Error: err}
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			return agentSavedMsg{Error: err}
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Charming-Key", "charming-secret-key")

		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return agentSavedMsg{Error: err}
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
			return agentSavedMsg{Error: fmt.Errorf("server error: %d", resp.StatusCode)}
		}

		return agentSavedMsg{}
	}
}

func (m rootModel) initDashboard() tea.Cmd {
	chat := tui.NewChatModel()
	chat.Model = m.wizard.Results.Model
	chat.APIKey = m.wizard.Results.APIKey
	
	switch {
	case strings.Contains(strings.ToLower(chat.Model), "gpt"):
		chat.Provider = "openai"
	case strings.Contains(strings.ToLower(chat.Model), "claude"):
		chat.Provider = "anthropic"
	case strings.Contains(strings.ToLower(chat.Model), "llama3"):
		chat.Provider = "ollama"
	default:
		chat.Provider = "openai"
	}

	chat.AddMessage(fmt.Sprintf("Agent %s initialized and saved!", m.wizard.Results.Name))
	
	win := tui.NewWindow(m.wizard.Results.Name, m.wizard.Results.Name, chat)
	win.Width = 60
	win.Height = 20
	win.X = 10
	win.Y = 5
	
	m.manager.AddWindow(win)
	m.manager.SetSize(m.width, m.height)
	
	return func() tea.Msg {
		return tea.WindowSizeMsg{Width: win.Width - 2, Height: win.Height - 2}
	}
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
