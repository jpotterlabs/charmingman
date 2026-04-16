package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"charmingman/backend/internal/provider"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHandleChat_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ps := provider.NewProviderService()
	h := NewChatHandler(ps)

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
	ps := provider.NewProviderService()
	h := NewChatHandler(ps)

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
