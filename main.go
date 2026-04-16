package main

import (
	"fmt"
	"os"

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
		state:  stateWizard,
		wizard: tui.NewWizardModel(),
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
			chat.AddMessage(fmt.Sprintf("Agent %s initialized with model %s", m.wizard.Results.Name, m.wizard.Results.Model))
			win := tui.NewWindow(m.wizard.Results.Name, m.wizard.Results.Name, chat)
			win.Width = 60
			win.Height = 20
			win.X = 10
			win.Y = 5
			m.manager.AddWindow(win)
			m.manager.SetSize(m.width, m.height)
			return m, nil
		}
		return m, cmd
	case stateDashboard:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "tab":
				if len(m.manager.Windows) > 0 {
					focusedIdx := -1
					for i, w := range m.manager.Windows {
						if w.Focused {
							focusedIdx = i
							break
						}
					}
					nextIdx := (focusedIdx + 1) % len(m.manager.Windows)
					m.manager.FocusWindow(m.manager.Windows[nextIdx].ID)
				}
			}
		case tea.MouseMsg:
			cmd := m.manager.HandleMouse(msg)
			return m, cmd
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

func (m rootModel) View() string {
	switch m.state {
	case stateWizard:
		return m.wizard.View()
	case stateDashboard:
		return m.manager.View()
	}
	return ""
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen(), tea.WithMouseAllMotion())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
