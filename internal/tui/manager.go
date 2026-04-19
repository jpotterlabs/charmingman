package tui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// Manager manages a collection of windows on an infinite canvas.
type Manager struct {
	Windows []*Window
	Width   int // Viewport width
	Height  int // Viewport height

	// Camera / Viewport offsets
	OffsetX int
	OffsetY int
	Zoom    float64 // 1.0 is normal

	LastMouseX int
	LastMouseY int
}

// NewManager creates a new window manager.
func NewManager() *Manager {
	return &Manager{
		Windows: make([]*Window, 0),
		Zoom:    1.0,
	}
}

// AddWindow adds a new window to the manager.
func (m *Manager) AddWindow(w *Window) {
	m.Windows = append(m.Windows, w)
}

// Update handles messages for all windows.
func (m *Manager) Update(msg tea.Msg) tea.Cmd {
	if wMsg, ok := msg.(WindowMsg); ok {
		for _, w := range m.Windows {
			if w.ID == wMsg.ID {
				return w.Update(wMsg.Msg)
			}
		}
		return nil
	}

	// Handle @mention routing
	if rMsg, ok := msg.(RouteMsg); ok {
		var cmds []tea.Cmd
		for _, w := range m.Windows {
			cmd := w.Update(rMsg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
		return tea.Batch(cmds...)
	}

	// Handle global panning/zooming keys (consume them)
	if msg, ok := msg.(tea.KeyPressMsg); ok {
		handled := false
		switch msg.String() {
		case "up":
			m.OffsetY -= 1
			handled = true
		case "down":
			m.OffsetY += 1
			handled = true
		case "left":
			m.OffsetX -= 2
			handled = true
		case "right":
			m.OffsetX += 2
			handled = true
		case "+":
			m.Zoom += 0.1
			handled = true
		case "-":
			if m.Zoom > 0.2 {
				m.Zoom -= 0.1
			}
			handled = true
		case "0":
			m.Zoom = 1.0
			m.OffsetX = 0
			m.OffsetY = 0
			handled = true
		}
		if handled {
			return nil
		}
	}

	var cmds []tea.Cmd
	for _, w := range m.Windows {
		cmd := w.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	return tea.Batch(cmds...)
}

// View renders all windows using a compositor and camera offsets.
func (m *Manager) View() tea.View {
	canvas := lipgloss.NewCanvas(m.Width, m.Height)
	comp := lipgloss.NewCompositor()

	for i, w := range m.Windows {
		// Calculate world-to-screen coordinates with scaling
		screenX := int(float64(w.X-m.OffsetX) * m.Zoom)
		screenY := int(float64(w.Y-m.OffsetY) * m.Zoom)
		scaledWidth := int(float64(w.Width) * m.Zoom)
		scaledHeight := int(float64(w.Height) * m.Zoom)

		// Basic culling using scaled dimensions
		if screenX+scaledWidth < 0 || screenX > m.Width || screenY+scaledHeight < 0 || screenY > m.Height {
			continue
		}

		layer := lipgloss.NewLayer(w.View().Content).
			ID(w.ID).
			X(screenX).
			Y(screenY).
			Z(i) 
		comp.AddLayers(layer)
	}

	return tea.NewView(canvas.Compose(comp).Render())
}

// HandleMouse handles mouse events for window management and canvas panning.
func (m *Manager) HandleMouse(msg tea.MouseMsg) tea.Cmd {
	mouse := msg.Mouse()
	dx := mouse.X - m.LastMouseX
	dy := mouse.Y - m.LastMouseY
	m.LastMouseX = mouse.X
	m.LastMouseY = mouse.Y

	switch msg.(type) {
	case tea.MouseClickMsg:
		if mouse.Button == tea.MouseLeft {
			// Check for clicks in reverse Z-order using scaled bounds
			for i := len(m.Windows) - 1; i >= 0; i-- {
				w := m.Windows[i]
				
				// Calculate scaled position on screen
				screenX := int(float64(w.X-m.OffsetX) * m.Zoom)
				screenY := int(float64(w.Y-m.OffsetY) * m.Zoom)
				scaledWidth := int(float64(w.Width) * m.Zoom)
				scaledHeight := int(float64(w.Height) * m.Zoom)

				// Resize handle (bottom right)
				if mouse.X >= screenX+scaledWidth-2 && mouse.X < screenX+scaledWidth && mouse.Y >= screenY+scaledHeight-1 && mouse.Y < screenY+scaledHeight {
					w.Resizing = true
					m.FocusWindow(w.ID)
					return nil
				}
				// Title bar (top row)
				if mouse.X >= screenX && mouse.X < screenX+scaledWidth && mouse.Y == screenY {
					w.Dragging = true
					m.FocusWindow(w.ID)
					return nil
				}
				// Anywhere else in window
				if mouse.X >= screenX && mouse.X < screenX+scaledWidth && mouse.Y >= screenY && mouse.Y < screenY+scaledHeight {
					m.FocusWindow(w.ID)
					return nil
				}
			}
		}
	case tea.MouseReleaseMsg:
		for _, w := range m.Windows {
			w.Dragging = false
			w.Resizing = false
		}
	case tea.MouseMotionMsg:
		anyDragging := false
		for _, w := range m.Windows {
			if w.Dragging {
				w.X += int(float64(dx) / m.Zoom)
				w.Y += int(float64(dy) / m.Zoom)
				anyDragging = true
				break
			}
			if w.Resizing {
				w.Width += int(float64(dx) / m.Zoom)
				w.Height += int(float64(dy) / m.Zoom)
				if w.Width < 10 { w.Width = 10 }
				if w.Height < 5 { w.Height = 5 }
				anyDragging = true
				break
			}
		}

		if !anyDragging && mouse.Button == tea.MouseLeft {
			// Background pan
			m.OffsetX -= int(float64(dx) / m.Zoom)
			m.OffsetY -= int(float64(dy) / m.Zoom)
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
			w.SetFocused(true)
		} else {
			w.SetFocused(false)
		}
	}
}

// SetSize sets the dimensions of the manager's workspace.
func (m *Manager) SetSize(width, height int) {
	m.Width = width
	m.Height = height
}
