package handler

import (
	"net/http"
	"strings"

	"charmingman/backend/internal/provider"
	"charm.land/fantasy"
	"github.com/gin-gonic/gin"
)

type ChatRequest struct {
	Provider string `json:"provider" binding:"required"`
	Model    string `json:"model" binding:"required"`
	Prompt   string `json:"prompt" binding:"required"`
}

type ChatResponse struct {
	Response string        `json:"response"`
	Usage    fantasy.Usage `json:"usage"`
	Error    string        `json:"error,omitempty"`
}

type ChatHandler struct {
	providerService *provider.ProviderService
}

func NewChatHandler(ps *provider.ProviderService) *ChatHandler {
	return &ChatHandler{
		providerService: ps,
	}
}

func (h *ChatHandler) HandleChat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.providerService.Chat(c.Request.Context(), req.Provider, req.Model, req.Prompt)
	if err != nil {
		if strings.Contains(err.Error(), "not registered") {
			c.JSON(http.StatusBadRequest, ChatResponse{Error: err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, ChatResponse{Error: err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, ChatResponse{
		Response: res.Content.Text(),
		Usage:    res.Usage,
	})
}
