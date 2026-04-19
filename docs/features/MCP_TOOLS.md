# MCP Tools & The Tool Belt

CharmingMan integrates with the **Model Context Protocol (MCP)**, an open standard for giving AI agents access to local tools and data sources. This allows your agents to perform actions like reading files, searching databases, or executing shell commands.

## 🛠️ How MCP Works in CharmingMan

1. **Server-Side**: You run one or more MCP servers (e.g., a filesystem server or a database server).
2. **Discovery**: The AI Gateway connects to these servers on startup and retrieves a list of available tools.
3. **Registration**: These tools are registered with the `fantasy` agent framework.
4. **Execution**: When an LLM decides to call a tool, the Gateway forwards the request to the local MCP server, executes the logic, and returns the result to the agent.

## ⚙️ Configuration

Configure your MCP servers in the `.env` file using the `MCP_SERVERS` variable. This should be a comma-separated list of commands to start each server.

**Example:**
```env
# Connect to a local filesystem and weather MCP server
MCP_SERVERS="npx @modelcontextprotocol/server-filesystem /path/to/docs, python -m weather_mcp"
```

## 🛠️ The Tool Belt (TUI)

CharmingMan includes a dedicated window called the **Tool Belt** for managing your active capabilities.

- **Discovery**: See all tools currently registered with your agents.
- **Descriptions**: View what each tool does and what parameters it expects.
- **Refresh**: Press `'r'` to re-scan the AI Gateway for updated tool definitions.

## 🚀 Future Roadmap: Authorization

*Coming soon:*
- **Per-Agent Allow-lists**: Restrict which agents can use which tools.
- **Human-in-the-Loop**: A popup in the TUI will ask for your confirmation before an agent executes sensitive tools (like shell commands or database writes).

## 🛠️ Implementation Details

- **Client**: `backend/internal/mcp/client.go` implements a lightweight JSON-RPC client over `stdin/stdout`.
- **Adapter**: `backend/internal/mcp/adapter.go` wraps MCP tools into the `fantasy.AgentTool` interface.
- **API**: `GET /api/v1/tools` exposes the list of registered tools to the TUI.
