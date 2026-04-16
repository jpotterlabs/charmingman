package tui

import (
	"slices"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// Manager manages a collection of windows.
type Manager struct {
	Windows []*Window
	Width   int
	Height  int
}

// NewManager creates a new window manager.
func NewManager() *Manager {
	return &Manager{
		Windows: make([]*Window, 0),
	}
}

// AddWindow adds a new window to the manager.
func (m *Manager) AddWindow(w *Window) {
	m.Windows = append(m.Windows, w)
}

// Update handles messages for all windows.
func (m *Manager) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	for _, w := range m.Windows {
		cmd := w.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	return tea.Batch(cmds...)
}

// View renders all windows using a compositor.
func (m *Manager) View() string {
	canvas := lipgloss.NewCanvas(m.Width, m.Height)
	comp := lipgloss.NewCompositor()

	for i, w := range m.Windows {
		layer := lipgloss.NewLayer(w.View()).
			ID(w.ID).
			X(w.X).
			Y(w.Y).
			Z(i) // Simple Z-index based on order
		comp.AddLayers(layer)
	}

	return canvas.Compose(comp).Render()
}

// HandleMouse handles mouse events for window management.
func (m *Manager) HandleMouse(msg tea.MouseMsg) tea.Cmd {
	switch msg.Type {
	case tea.MouseLeft:
		// Check for clicks in reverse Z-order
		for i := len(m.Windows) - 1; i >= 0; i-- {
			w := m.Windows[i]
			if w.IsInResizeHandle(msg.X, msg.Y) {
				w.Resizing = true
				m.FocusWindow(w.ID)
				return nil
			}
			if w.IsInTitleBar(msg.X, msg.Y) {
				w.Dragging = true
				m.FocusWindow(w.ID)
				return nil
			}
			// Check if click is anywhere in the window to focus
			if msg.X >= w.X && msg.X < w.X+w.Width && msg.Y >= w.Y && msg.Y < w.Y+w.Height {
				m.FocusWindow(w.ID)
				return nil
			}
		}
	case tea.MouseRelease:
		for _, w := range m.Windows {
			w.Dragging = false
			w.Resizing = false
		}
	case tea.MouseMotion:
		for _, w := range m.Windows {
			if w.Dragging {
				w.X += msg.DX
				w.Y += msg.DY
				return nil
			}
			if w.Resizing {
				w.Width += msg.DX
				w.Height += msg.DY
				if w.Width < 10 {
					w.Width = 10
				}
				if w.Height < 5 {
					w.Height = 5
				}
				return nil
			}
		}
	}
	return nil
}

// FocusWindow brings a window to the front.
func (m *Manager) FocusWindow(id string) {
	for i, w := range m.Windows {
		if w.ID == id {
			// Move to the end of the slice (top Z-index)
			m.Windows = append(m.Windows[:i], m.Windows[i+1:]...)
			m.Windows = append(m.Windows, w)
			w.Focused = true
		} else {
			w.Focused = false
		}
	}
}

// SetSize sets the dimensions of the manager's workspace.
func (m *Manager) SetSize(width, height int) {
	m.Width = width
	m.Height = height
}
