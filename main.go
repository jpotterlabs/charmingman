package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"charmingman/internal/config"
	"charmingman/internal/tui"
	tea "charm.land/bubbletea/v2"
	"github.com/google/uuid"
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
	roomID  string
}

func initialModel() rootModel {
	return rootModel{
		state:   stateWizard,
		wizard:  tui.NewWizardModel(),
		manager: tui.NewManager(),
		roomID:  uuid.New().String(), // Generate a shared room ID for the session
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
			return m, m.resizeDashboard()
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
			return m, m.initDashboardFromLayout()
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

		baseURL := strings.TrimRight(rawURL, "/")
		baseURL = strings.TrimSuffix(baseURL, "/api/v1/chat")
		baseURL = strings.TrimSuffix(baseURL, "/api/v1/agents")
		baseURL = strings.TrimSuffix(baseURL, "/chat")
		baseURL = strings.TrimSuffix(baseURL, "/agents")
		
		finalURL := baseURL + "/api/v1/agents"
		if !strings.Contains(finalURL, "://") {
			finalURL = "http://" + finalURL
		}

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

func (m *rootModel) calculateGeometry() (chatX, chatY, chatW, chatH, docX, docY, docW, docH int) {
	if m.width >= 150 {
		chatW = 60
		chatH = m.height - 4
		chatX = 2
		chatY = 1

		docW = m.width - chatW - 6
		docH = m.height - 4
		docX = chatW + 4
		docY = 1
	} else {
		chatW = m.width - 4
		chatH = (m.height / 2) - 2
		chatX = 2
		chatY = 1

		docW = m.width - 4
		docH = (m.height / 2) - 2
		docX = 2
		docY = chatH + 2
	}

	if chatW < 20 { chatW = 20 }
	if chatH < 10 { chatH = 10 }
	if docW < 20 { docW = 20 }
	if docH < 10 { docH = 10 }
	
	return
}

func (m rootModel) initDashboardFromLayout() tea.Cmd {
	layout, err := config.LoadLayout("layout.yaml")
	if err != nil {
		log.Printf("Failed to load layout.yaml: %v. Falling back to default.", err)
		return m.initDashboardDefault()
	}

	var cmds []tea.Cmd
	for _, winCfg := range layout.Windows {
		var winModel tea.Model
		switch winCfg.Type {
		case "chat":
			chat := tui.NewChatModel(winCfg.ID)
			chat.Provider = winCfg.Config.Provider
			chat.Model = winCfg.Config.Model
			chat.UseRAG = winCfg.Config.UseRAG
			chat.RoomID = m.roomID // Shared room ID
			
			if winCfg.ID == "primary-chat" {
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
				}
				chat.AddMessage(fmt.Sprintf("Agent %s initialized from Wizard", m.wizard.Results.Name))
			} else {
				chat.AddMessage(fmt.Sprintf("Agent %s initialized from Layout", winCfg.Title))
			}
			winModel = chat

		case "document":
			doc := tui.NewDocumentModel(winCfg.Config.Content)
			winModel = doc

		default:
			continue
		}

		win := tui.NewWindow(winCfg.ID, winCfg.Title, winModel)
		win.X = winCfg.X
		win.Y = winCfg.Y
		win.Width = winCfg.Width
		win.Height = winCfg.Height
		
		m.manager.AddWindow(win)
		
		cmds = append(cmds, func(id string, w, h int) func() tea.Msg {
			return func() tea.Msg {
				return tui.WindowMsg{ID: id, Msg: tea.WindowSizeMsg{Width: w - 2, Height: h - 2}}
			}
		}(winCfg.ID, win.Width, win.Height))
	}

	return tea.Batch(cmds...)
}

func (m rootModel) initDashboardDefault() tea.Cmd {
	chat := tui.NewChatModel("primary-chat")
	chat.Model = m.wizard.Results.Model
	chat.APIKey = m.wizard.Results.APIKey
	chat.UseRAG = m.wizard.Results.UseRAG
	chat.RoomID = m.roomID
	
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

	chat.AddMessage(fmt.Sprintf("Agent %s initialized (Default Layout)", m.wizard.Results.Name))
	
	chatWin := tui.NewWindow("primary-chat", m.wizard.Results.Name, chat)
	chatWin.X = 2
	chatWin.Y = 1
	chatWin.Width = 60
	chatWin.Height = 20
	
	m.manager.AddWindow(chatWin)
	
	return func() tea.Msg {
		return tui.WindowMsg{ID: "primary-chat", Msg: tea.WindowSizeMsg{Width: chatWin.Width - 2, Height: chatWin.Height - 2}}
	}
}

func (m rootModel) resizeDashboard() tea.Cmd {
	var cmds []tea.Cmd
	for _, w := range m.manager.Windows {
		if w.X + w.Width > m.width {
			if w.Width > m.width - 2 {
				w.Width = m.width - 2
			}
			if w.X + w.Width > m.width {
				w.X = m.width - w.Width - 1
			}
		}
		if w.Y + w.Height > m.height {
			if w.Height > m.height - 2 {
				w.Height = m.height - 2
			}
			if w.Y + w.Height > m.height {
				w.Y = m.height - w.Height - 1
			}
		}
		
		cmds = append(cmds, func(id string, width, height int) func() tea.Msg {
			return func() tea.Msg {
				return tui.WindowMsg{ID: id, Msg: tea.WindowSizeMsg{Width: width - 2, Height: height - 2}}
			}
		}(w.ID, w.Width, w.Height))
	}
	return tea.Batch(cmds...)
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
