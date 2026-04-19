package mcp

import (
	"context"
	"encoding/json"

	"charm.land/fantasy"
)

// MCPToolWrapper wraps an MCP tool to implement fantasy.AgentTool.
type MCPToolWrapper struct {
	client *Client
	tool   Tool
}

func NewMCPToolWrapper(client *Client, tool Tool) *MCPToolWrapper {
	return &MCPToolWrapper{
		client: client,
		tool:   tool,
	}
}

func (w *MCPToolWrapper) Info() fantasy.ToolInfo {
	var params map[string]any
	json.Unmarshal(w.tool.InputSchema, &params)
	
	return fantasy.ToolInfo{
		Name:        w.tool.Name,
		Description: w.tool.Description,
		Parameters:  params,
	}
}

func (w *MCPToolWrapper) Run(ctx context.Context, params fantasy.ToolCall) (fantasy.ToolResponse, error) {
	// Parse input from fantasy params. Input is expected to be a JSON string.
	var input map[string]interface{}
	if err := json.Unmarshal([]byte(params.Input), &input); err != nil {
		return fantasy.NewTextErrorResponse(err.Error()), nil
	}

	callParams := struct {
		Name      string                 `json:"name"`
		Arguments map[string]interface{} `json:"arguments"`
	}{
		Name:      w.tool.Name,
		Arguments: input,
	}

	resp, err := w.client.Call(ctx, "tools/call", callParams)
	if err != nil {
		return fantasy.NewTextErrorResponse(err.Error()), nil
	}

	var result struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return fantasy.NewTextErrorResponse(err.Error()), nil
	}

	if len(result.Content) > 0 {
		return fantasy.NewTextResponse(result.Content[0].Text), nil
	}

	return fantasy.NewTextResponse("No content returned from tool"), nil
}

func (w *MCPToolWrapper) ProviderOptions() fantasy.ProviderOptions {
	return nil
}

func (w *MCPToolWrapper) SetProviderOptions(opts fantasy.ProviderOptions) {
}
