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
	stateSaveFailed
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
	err     error
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
	// Capture window size updates regardless of state
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = msg.Width
		m.height = msg.Height
		if m.state == stateDashboard {
			m.manager.SetSize(msg.Width, msg.Height)
		}
	}

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
				m.err = msg.Error
				m.state = stateSaveFailed
				return m, nil
			}
			m.state = stateDashboard
			m.manager.SetSize(m.width, m.height)
			return m, m.initDashboard()
		}
		return m, nil

	case stateSaveFailed:
		if msg, ok := msg.(tea.KeyPressMsg); ok {
			if msg.String() == "enter" {
				m.state = stateSavingAgent
				return m, m.saveAgent()
			}
			if msg.String() == "q" || msg.String() == "ctrl+c" {
				return m, tea.Quit
			}
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
	case stateSaveFailed:
		v.SetContent(fmt.Sprintf("Error saving agent: %v\n\nPress Enter to retry, or q to quit.", m.err))
	case stateDashboard:
		v.SetContent(m.manager.View().Content)
	}
	return v
}

func (m rootModel) saveAgent() tea.Cmd {
	return func() tea.Msg {
		rawURL := os.Getenv("GATEWAY_URL")
		if rawURL == "" {
			rawURL = "http://localhost:8090"
		}

		// Normalize URL to /api/v1/agents
		baseURL := strings.TrimRight(rawURL, "/")
		baseURL = strings.TrimSuffix(baseURL, "/api/v1/chat")
		baseURL = strings.TrimSuffix(baseURL, "/api/v1/agents")
		baseURL = strings.TrimSuffix(baseURL, "/chat")
		baseURL = strings.TrimSuffix(baseURL, "/agents")
		
		finalURL := baseURL + "/api/v1/agents"
		if !strings.Contains(finalURL, "://") {
			finalURL = "http://" + finalURL
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

		payload := map[string]interface{}{
			"name":     m.wizard.Results.Name,
			"model":    model,
			"provider": provider,
			"persona":  m.wizard.Results.Persona,
			"api_key":  m.wizard.Results.APIKey,
			"use_rag":  m.wizard.Results.UseRAG,
		}

		jsonData, err := json.Marshal(payload)
		if err != nil {
			return agentSavedMsg{Error: err}
		}

		req, err := http.NewRequest("POST", finalURL, bytes.NewBuffer(jsonData))
		if err != nil {
			return agentSavedMsg{Error: err}
		}
		req.Header.Set("Content-Type", "application/json")
		
		apiKey := os.Getenv("GATEWAY_API_KEY")
		if apiKey == "" {
			return agentSavedMsg{Error: fmt.Errorf("GATEWAY_API_KEY environment variable is not set")}
		}
		req.Header.Set("X-Charming-Key", apiKey)

		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return agentSavedMsg{Error: err}
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
			var errResp struct {
				Error string `json:"error"`
			}
			json.NewDecoder(resp.Body).Decode(&errResp)
			return agentSavedMsg{Error: fmt.Errorf("server error (%d): %s", resp.StatusCode, errResp.Error)}
		}

		return agentSavedMsg{}
	}
}

func (m rootModel) initDashboard() tea.Cmd {
	// 1. Chat Window
	chat := tui.NewChatModel()
	chat.Model = m.wizard.Results.Model
	chat.APIKey = m.wizard.Results.APIKey
	chat.UseRAG = m.wizard.Results.UseRAG
	
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
	
	chatWin := tui.NewWindow(m.wizard.Results.Name, m.wizard.Results.Name, chat)
	chatWin.Width = 60
	chatWin.Height = 20
	chatWin.X = 5
	chatWin.Y = 2
	
	m.manager.AddWindow(chatWin)

	// 2. Document Preview Window (The Stage)
	doc := tui.NewDocumentModel("# Knowledge Base\n\nWelcome to your agent's knowledge base. You can upload documents here to ground your assistant in specific context.")
	docWin := tui.NewWindow("stage", "The Stage", doc)
	docWin.Width = 80
	docWin.Height = 25
	docWin.X = 70
	docWin.Y = 5
	
	m.manager.AddWindow(docWin)
	
	// Create individual WindowSizeMsg commands for each window's sub-model
	return tea.Batch(
		func() tea.Msg {
			// Calculate content area for chatWin
			return tui.WindowMsg{
				ID: chatWin.ID,
				Msg: tea.WindowSizeMsg{
					Width:  chatWin.Width - 2,
					Height: chatWin.Height - 2,
				},
			}
		},
		func() tea.Msg {
			// Calculate content area for docWin
			return tui.WindowMsg{
				ID: docWin.ID,
				Msg: tea.WindowSizeMsg{
					Width:  docWin.Width - 2,
					Height: docWin.Height - 2,
				},
			}
		},
	)
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
