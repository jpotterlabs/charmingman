Build a configurable "Window Manager" using **Bubble Tea** and **Lipgloss**. Make use of comprehensive mouse support.

Build out 5 pre-set user interfaceds for CharmingManusing the abbove Window Manager.

**Requirements**
1.  **Mouse Support:** You have to explicitly tell your program to listen for mouse events by setting the `MouseMode` in your `View()` (e.g., `tea.MouseModeAllMotion`).
2.  **Tracking State:** You define a "window" as a struct with an `X`, `Y`, `Width`, and `Height`.
3.  **Handling Events:**
    * **Moving:** When a `tea.MouseClickMsg` hits your "title bar" area, you set a `dragging` flag to true. As long as the button is held, you use `tea.MouseMotionMsg` to update the window's `X` and `Y` coordinates based on the cursor's movement.
    * **Resizing:** Similar to moving, but you listen for clicks on the border or a corner "handle." When the mouse moves, you increment/decrement the `Width` and `Height` values of your component.
4.  **Rendering:** You use **Lipgloss** to render the window at those specific coordinates using `lipgloss.Place()` or by calculating the padding/margins dynamically.
5. Pre-build components assemble to become pre-set User Interfaces.
6. Implement "double-click to maximize" a pane

**Watch this space**

* **`Lipgloss v2`:** Recently introduced a `Canvas` and `Layer` system that makes "Z-index" (stacking windows on top of each other) and coordinate-based positioning much easier than the original version.
* **The "Bubbles" Viewport:** While primarily for scrolling, it’s often used as the "inner pane" of a custom window.

**Pro-tip:** Most TUI users still prefer keyboard shortcuts, so if you build mouse-resizable windows, definitely map `Ctrl+Arrows` or similar to do the same thing for the "no-mouse" crowd!


That is a massive leap forward for the **CharmingMan** ecosystem. By treating the UI as a **modular, shareable asset**, you’re moving from a "Chat App" to a "TUI Framework for Agents."

In the Go/Charm world, this essentially means creating a **Spec-Driven UI** where the layout isn't hardcoded, but rather interpreted from a configuration file (like YAML or JSON) or even a Go plugin.

---

### 1. The "CharmingManifest" (UI-as-Code)
To make UIs shareable, you need a standardized way to describe them. Users could define a `layout.yaml` that CharmingMan parses to build the `Bubble Tea` view:

```yaml
name: "The Researcher"
author: "RomanHistorian90"
layout:
  type: "split-horizontal"
  panes:
    - id: "main-chat"
      width: 70%
      component: "viewport"
    - id: "sidebar"
      width: 30%
      type: "stack-vertical"
      components:
        - "token-monitor"
        - "agent-status"
        - "mcp-explorer"
```

### 2. A "Pluggable" Component Gallery
If users can create their own UIs, you’ll want a library of "widgets" they can drop in. Since you’re already using **Bubbles** and **Huh?**, these widgets are already modular:
* **The Log-Streamer:** A terminal-style tail of agent reasoning.
* **The Markdown-Stage:** A `Glow` powered area for the final output.
* **The Tool-Belt:** A grid of buttons (clickable via mouse) to trigger specific MCP tools or local Python scripts via `uv`.

### 3. Sharing via "Charming Registry"
Imagine a command like `charmingman ui install big-agi-clone`. 
* **Implementation:** The app fetches the manifest from a GitHub repo or a central registry.
* **Theming:** Using **Catwalk**, the shared UI automatically adapts to the user's local terminal colors, so a UI designed in "Dracula" still looks great for a user on "Solarized Light."

### 4. The "Live Designer" Mode
Since you mentioned mouse support (resizing/moving), you could have a **"Design Mode"** within the app:
1.  **Enter Design Mode:** `Ctrl+D`.
2.  **Interact:** Drag borders to resize, right-click to change a pane's component (e.g., swap a Chat window for a PDF viewer).
3.  **Export:** `Ctrl+S` saves the current state as a new manifest file that can be shared.

---

### Technical Challenge: Dynamic Typing in Go
Since Go is a statically typed language, creating a truly "dynamic" UI that loads at runtime can be tricky. You have two main paths:
* **The Interpretive Path:** You write a "Master Model" in `Bubble Tea` that knows how to render a variety of components based on the YAML spec.
* **The Plugin Path:** You use Go's `plugin` package or **WebAssembly (Wasm)** to allow users to compile their own UI components and "hot-swap" them into the running CharmingMan binary.

### The "Big-AGI" Connection
Just as **big-AGI** allows users to customize their "Models," "Personas," and "Prompts," **CharmingMan** would allow them to customize the **"Glass"** through which they interact with their agents. This is perfect for your **Spec Driven Design** tool—one user might want a UI optimized for writing code, while another wants one optimized for historical research.

Does the idea of a **YAML-based layout engine** sound like the right level of complexity, or were you thinking of something even more dynamic?
