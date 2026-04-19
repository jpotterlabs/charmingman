# Infinity Canvas & Spatial UI

The **Infinity Canvas** is CharmingMan's spatial workspace. Instead of a traditional linear or tabbed interface, CharmingMan provides a multi-directional coordinate plane where you can arrange agents, documents, and tools as floating windows.

## 🌌 Core Mechanics

### 1. Spatial Coordinate System
Every window in the TUI has a "World Position" defined by absolute `(X, Y)` coordinates.
- **World Coordinates**: High-resolution positions on the infinite plane.
- **Screen Coordinates**: The actual grid cells on your terminal window.

### 2. Camera System
The `Manager` acts as a camera, maintaining a viewport over the world:
- **OffsetX / OffsetY**: These values determine which part of the world is currently visible.
- **Zoom Level**: Scales the world to fit more or fewer components on screen.

### 3. Coordinate Mapping
The mapping from world to screen is calculated dynamically:
```go
screenX := int(float64(window.X - OffsetX) * Zoom)
screenY := int(float64(window.Y - OffsetY) * Zoom)
```
This ensures that as you pan or zoom, the spatial relationships between your windows remain consistent.

## 🕹️ Interaction & Navigation

CharmingMan provides intuitive controls for exploring the canvas:

| Action | Control |
|--------|---------|
| **Pan Camera** | Arrow Keys or Drag Background (Left Mouse) |
| **Zoom In/Out** | `+` / `-` keys |
| **Reset View** | `0` key (Home) |
| **Move Window**| Drag Window Title Bar (Left Mouse) |
| **Resize Window**| Drag Bottom-Right Corner (Left Mouse) |
| **Cycle Focus**| `Tab` key |

### Z-Index Stacking
Windows support layering. When you click on a window or focus it via `Tab`, the `Manager` automatically brings it to the front of the stack, ensuring your active workspace is always visible.

## 📝 YAML Layouts (`layout.yaml`)

You can define your ideal workspace configuration in `layout.yaml`. CharmingMan parses this file on startup to restore your preferred arrangement of agents and tools.

**Example Schema:**
```yaml
windows:
  - id: "primary-chat"
    title: "Project Lead"
    type: "chat"
    x: 2
    y: 1
    width: 60
    height: 25
    focused: true
```

### Auto-Rescaling
CharmingMan includes logic to ensure that windows defined in `layout.yaml` are visible on your current terminal. If a window's dimensions exceed your terminal size at startup, the loader will automatically scale it down to fit the available bounds.

## 🚀 Semantic Zoom
*Coming soon in Phase 4 iterations.*
- **High Zoom**: Detailed view of chat history and input.
- **Mid Zoom**: Metadata summary and recent message preview.
- **Low Zoom**: Icon-only view for high-density "God View" of the swarm.
