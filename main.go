package main

import (
	"fmt"
	"os"
	"strings"

	"charmingman/internal/tui"
	tea "charm.land/bubbletea/v2"
)

type state int

const (
	stateWizard state = iota
	stateDashboard
)

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
			m.state = stateDashboard
			// Initialize dashboard with agent
			chat := tui.NewChatModel()
			
			// Map model names to providers for the backend
			chat.Model = m.wizard.Results.Model
			chat.APIKey = m.wizard.Results.APIKey
			
			switch {
			case strings.Contains(strings.ToLower(chat.Model), "gpt"):
				chat.Provider = "openai"
			case strings.Contains(strings.ToLower(chat.Model), "claude"):
				chat.Provider = "anthropic"
			case strings.Contains(strings.ToLower(chat.Model), "llama3"):
				chat.Provider = "ollama" // Default local to ollama for now
			default:
				chat.Provider = "openai"
			}

			chat.AddMessage(fmt.Sprintf("Agent %s initialized with model %s on provider %s", m.wizard.Results.Name, chat.Model, chat.Provider))
			
			win := tui.NewWindow(m.wizard.Results.Name, m.wizard.Results.Name, chat)
			win.Width = 60
			win.Height = 20
			win.X = 10
			win.Y = 5
			
			m.manager.AddWindow(win)
			m.manager.SetSize(m.width, m.height)
			
			// Manually send initial size to the new chat model
			chat.Update(tea.WindowSizeMsg{Width: win.Width - 2, Height: win.Height - 2})
			
			return m, nil
		}
		return m, cmd
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
	case stateDashboard:
		v.SetContent(m.manager.View().Content)
	}
	return v
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
