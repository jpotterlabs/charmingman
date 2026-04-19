package tui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type tool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type toolsMsg []tool

type ToolBeltModel struct {
	tools  []tool
	err    error
	ready  bool
	width  int
	height int
}

func NewToolBeltModel() *ToolBeltModel {
	return &ToolBeltModel{}
}

func (m *ToolBeltModel) Init() tea.Cmd {
	return m.fetchTools
}

func (m *ToolBeltModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true

	case toolsMsg:
		m.tools = msg
		return m, nil

	case tea.KeyPressMsg:
		if msg.String() == "r" {
			return m, m.fetchTools
		}
	}
	return m, nil
}

func (m *ToolBeltModel) View() tea.View {
	if !m.ready {
		return tea.NewView("Loading tools...")
	}

	if m.err != nil {
		return tea.NewView(fmt.Sprintf("Error: %v", m.err))
	}

	if len(m.tools) == 0 {
		return tea.NewView("No active MCP tools found.\n\nPress 'r' to refresh.")
	}

	style := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, true, false).
		Padding(0, 1)

	var content string
	content += "Active MCP Tools:\n\n"
	for _, t := range m.tools {
		content += style.Render(fmt.Sprintf("🛠️ %s: %s", t.Name, t.Description)) + "\n"
	}
	content += "\nPress 'r' to refresh."

	return tea.NewView(content)
}

func (m *ToolBeltModel) fetchTools() tea.Msg {
	url := os.Getenv("GATEWAY_URL")
	if url == "" {
		url = "http://localhost:8090/api/v1/tools"
	} else {
		url = strings.TrimSuffix(url, "/chat") + "/tools"
	}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Charming-Key", os.Getenv("GATEWAY_API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch tools: %w", err)
	}
	defer resp.Body.Close()

	var tools toolsMsg
	json.NewDecoder(resp.Body).Decode(&tools)
	return tools
}
