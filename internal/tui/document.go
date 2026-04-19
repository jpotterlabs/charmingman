package tui

import (
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
)

type DocumentModel struct {
	viewport viewport.Model
	content  string
	ready    bool
	width    int
	height   int
}

func NewDocumentModel(content string) *DocumentModel {
	return &DocumentModel{
		content: content,
	}
}

func (m *DocumentModel) Init() tea.Cmd {
	return nil
}

func (m *DocumentModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		if !m.ready {
			m.viewport = viewport.New(viewport.WithWidth(m.width), viewport.WithHeight(m.height))
			m.viewport.SetContent(m.content)
			m.ready = true
		} else {
			m.viewport.SetWidth(m.width)
			m.viewport.SetHeight(m.height)
		}
	}

	if m.ready {
		m.viewport, cmd = m.viewport.Update(msg)
	}
	return m, cmd
}

func (m *DocumentModel) View() tea.View {
	if !m.ready {
		return tea.NewView("Loading document...")
	}
	return tea.NewView(m.viewport.View())
}

func (m *DocumentModel) SetContent(content string) {
	m.content = content
	if m.ready {
		m.viewport.SetContent(m.content)
	}
}
