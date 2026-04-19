# CharmingMan User Interfaces & Layouts

CharmingMan features a modular, configurable "Window Manager" built with **Bubble Tea** and **Lipgloss**. This system supports both a traditional tiling window manager and a revolutionary **Infinity Canvas**.

## 1. The Infinity Canvas

The Infinity Canvas is a non-linear spatial environment. It moves beyond 1D streams (standard chat) and 2D grids (dashboards) to provide a **Spatial Workspace**.

### Key Concepts:
- **Spatial Layout**: Windows are placed at absolute $(x, y)$ coordinates in a global "world."
- **Camera Panning**: Move the view using `Ctrl + Arrow Keys` or mouse drag on the workspace background.
- **Camera Zooming**: Scale the entire workspace using `Ctrl + +` and `Ctrl + -` (or mouse wheel).
- **World-to-Screen Mapping**: The engine dynamically translates global coordinates and zoom factors into terminal screen characters.

## 2. YAML-Driven Layouts (`layout.yaml`)

To make UIs shareable and maintainable, CharmingMan uses a standardized `layout.yaml` manifest.

### Schema Overview:

```yaml
name: "The Researcher"
author: "CharmingTeam"
canvas:
  background: "dots" # Background pattern (dots, grid, none)
  zoom: 1.0
  camera:
    x: 0
    y: 0
panes:
  - id: "main-chat"
    component: "chat"
    x: 10
    y: 5
    width: 60
    height: 20
    fixed: false # Can be dragged
    props:
      agent_id: "researcher-v1"
  - id: "stage"
    component: "stage"
    x: 75
    y: 5
    width: 40
    height: 30
    fixed: true
```

### Advanced Features:

1. **Semantic Validation**:
   The loader performs checks to ensure:
   - All `component` types are registered (e.g., `chat`, `stage`, `status`).
   - `id`s are unique.
   - Initial positions are within reasonable bounds.

2. **Auto-Rescaling**:
   When the terminal window size changes, CharmingMan calculates a scale factor. Window dimensions ($w, h$) and coordinates ($x, y$) are rescaled to ensure the layout remains visually consistent even on smaller screens.

3. **Inter-Component Messaging**:
   Components can communicate via a shared message bus. For example, a "mention" in a `chat` component can trigger a focus change in another agent's window.

## 3. Mouse Support & Interactions

CharmingMan provides high-fidelity mouse integration:
- **Dragging**: Click and hold a window's title bar to move it across the canvas.
- **Resizing**: Click and drag the borders or corner handles of a window.
- **Focusing**: Clicking anywhere inside a window brings it to the front (Z-index update).
- **Maximizing**: Double-click a title bar to toggle full-screen focus.

## 4. Keyboard Navigation

For "no-mouse" enthusiasts, keyboard shortcuts provide full control:
- `Tab`: Cycle focus between windows.
- `Ctrl + N`: Create a new agent window.
- `Ctrl + Arrow Keys`: Pan the canvas camera.
- `Ctrl + W`: Close the active window.
- `/`: Enter "Command Mode" for model switching or tool execution.

## 5. Visual Styling with `Catwalk` & `Lipgloss`

CharmingMan uses `Catwalk` to ensure that agent themes (Mauve for GPT-4, Green for Ollama) look great on any terminal color scheme (Dracula, Nord, Solarized). Borders, padding, and accent colors are all dynamically calculated based on the active terminal theme.
