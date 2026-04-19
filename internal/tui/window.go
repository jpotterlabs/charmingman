package tui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Focusable interface {
	SetFocused(bool)
}

// WindowMsg is a message routed to a specific window.
type WindowMsg struct {
	ID  string
	Msg tea.Msg
}

// Window represents a single window in the TUI.
type Window struct {
	ID      string
	Title   string
	X, Y    int
	Width   int
	Height  int
	Focused bool
	Model   tea.Model
	Style   lipgloss.Style

	Dragging bool
	Resizing bool
}

func (w *Window) IsInTitleBar(x, y int) bool {
	return x >= w.X && x < w.X+w.Width && y == w.Y
}

func (w *Window) IsInResizeHandle(x, y int) bool {
	return x >= w.X+w.Width-2 && x < w.X+w.Width && y >= w.Y+w.Height-1 && y < w.Y+w.Height
}

// NewWindow creates a new window with the given ID and model.
func NewWindow(id, title string, model tea.Model) *Window {
	return &Window{
		ID:    id,
		Title: title,
		Model: model,
		Style: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(0, 1),
	}
}

func (w *Window) SetFocused(f bool) {
	w.Focused = f
	if focusable, ok := w.Model.(Focusable); ok {
		focusable.SetFocused(f)
	}
}

// Update updates the window's model.
func (w *Window) Update(msg tea.Msg) (tea.Cmd) {
	var cmd tea.Cmd
	w.Model, cmd = w.Model.Update(msg)
	return cmd
}

// View renders the window.
func (w *Window) View() tea.View {
	style := w.Style
	if w.Focused {
		style = style.BorderForeground(lipgloss.Color("#EE6FF8"))
	}

	content := w.Model.View().Content
	rendered := style.
		Width(w.Width - style.GetHorizontalFrameSize()).
		Height(w.Height - style.GetVerticalFrameSize()).
		Render(content)
	
	return tea.NewView(rendered)
}
