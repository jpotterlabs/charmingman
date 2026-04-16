package tui

import (
	"charm.land/huh/v2"
	tea "charm.land/bubbletea/v2"
)

type WizardModel struct {
	form    *huh.Form
	done    bool
	Results AgentConfig
}

type AgentConfig struct {
	Name      string
	Model     string
	Persona   string
	APIKey    string
}

func NewWizardModel() *WizardModel {
	config := AgentConfig{}
	
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Agent Name").
				Value(&config.Name),
			huh.NewSelect[string]().
				Title("Model").
				Options(
					huh.NewOption("GPT-4o", "gpt-4o"),
					huh.NewOption("Claude 3.5 Sonnet", "claude-3-5-sonnet"),
					huh.NewOption("Llama 3 (Local)", "llama3"),
				).
				Value(&config.Model),
			huh.NewText().
				Title("Persona").
				Value(&config.Persona),
			huh.NewInput().
				Title("API Key").
				EchoMode(huh.EchoPassword).
				Value(&config.APIKey),
		),
	)

	return &WizardModel{
		form:    form,
		Results: config,
	}
}

func (m *WizardModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m *WizardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	if m.form.State == huh.StateCompleted {
		m.done = true
	}

	return m, cmd
}

func (m *WizardModel) View() string {
	return m.form.View()
}

func (m *WizardModel) IsDone() bool {
	return m.done
}
