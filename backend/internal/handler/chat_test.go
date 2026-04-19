package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"charmingman/backend/internal/document"
	"charmingman/backend/internal/provider"
	"charmingman/backend/internal/vector"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleChat_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ps := provider.NewProviderService(nil)
	h := NewChatHandler(ps, nil, nil)

	r := gin.Default()
	r.POST("/chat", h.HandleChat)

	reqBody := []byte(`{"invalid": "json"}`)
	req, _ := http.NewRequest("POST", "/chat", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleChat_ProviderNotRegistered(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ps := provider.NewProviderService(nil)
	h := NewChatHandler(ps, nil, nil)

	r := gin.Default()
	r.POST("/chat", h.HandleChat)

	reqBody, _ := json.Marshal(ChatRequest{
		Provider: "nonexistent",
		Model:    "gpt-4",
		Prompt:   "Hello",
	})
	req, _ := http.NewRequest("POST", "/chat", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp ChatResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Contains(t, resp.Error, "provider nonexistent not registered")
}

// MockEmbedder for testing
type mockEmbedder struct{}

func (m *mockEmbedder) Embed(ctx context.Context, text string) ([]float32, error) {
	return []float32{1.0, 0.0, 0.0}, nil
}
func (m *mockEmbedder) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
	return [][]float32{{1.0, 0.0, 0.0}}, nil
}

func TestHandleChat_RAGContextInjection(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// Create a mock openaicompat server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var reqBody map[string]interface{}
		json.Unmarshal(body, &reqBody)
		
		prompt := ""
		// Look for messages in the OpenAI-compatible request
		if msgs, ok := reqBody["messages"].([]interface{}); ok && len(msgs) > 0 {
			if lastMsg, ok := msgs[len(msgs)-1].(map[string]interface{}); ok {
				prompt, _ = lastMsg["content"].(string)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id": "chatcmpl-123",
			"object": "chat.completion",
			"created": 1677652288,
			"model": "test-model",
			"choices": []map[string]interface{}{
				{
					"index": 0,
					"message": map[string]interface{}{
						"role": "assistant",
						"content": "ECHO: " + prompt,
					},
					"finish_reason": "stop",
				},
			},
			"usage": map[string]interface{}{
				"prompt_tokens": 9,
				"completion_tokens": 12,
				"total_tokens": 21,
			},
		})
	}))
	defer mockServer.Close()

	// 1. Setup Provider using RegisterLocal
	ps := provider.NewProviderService(nil)
	err := ps.RegisterLocal("mock", mockServer.URL)
	assert.NoError(t, err)

	// 2. Setup Document Service with LocalStore
	store := vector.NewLocalStore()
	store.Add(context.Background(), []vector.Vector{
		{
			ID: "doc1", 
			Values: []float32{1.0, 0.0, 0.0}, 
			Metadata: map[string]interface{}{
				"content": "This is relevant RAG context.",
				"document_id": "123",
			},
		},
	})
	
	ds := document.NewDocumentService(nil, &mockEmbedder{}, store, "")
	
	h := NewChatHandler(ps, ds, nil)

	r := gin.Default()
	r.POST("/chat", h.HandleChat)

	// 3. Test with RAG enabled
	reqData, _ := json.Marshal(ChatRequest{
		Provider: "mock",
		Model:    "test-model",
		Prompt:   "What is this?",
		UseRAG:   true,
	})
	req, _ := http.NewRequest("POST", "/chat", bytes.NewBuffer(reqData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var resp ChatResponse
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	// Check if RAG sources are returned
	require.Len(t, resp.Sources, 1)
	require.Equal(t, "This is relevant RAG context.", resp.Sources[0].Content)

	// Check if prompt was injected
	assert.Contains(t, resp.Response, "Use the following pieces of context to answer")
	assert.Contains(t, resp.Response, "This is relevant RAG context.")
	assert.Contains(t, resp.Response, "Question: What is this?")
}