package handler

import (
	"net/http"
	"strings"

	"charmingman/backend/internal/document"
	"charmingman/backend/internal/provider"
	"charm.land/fantasy"
	"github.com/gin-gonic/gin"
)

type ChatRequest struct {
	Provider string `json:"provider" binding:"required"`
	Model    string `json:"model" binding:"required"`
	Prompt   string `json:"prompt" binding:"required"`
	UseRAG   bool   `json:"use_rag"`
}

type ChatResponse struct {
	Response string        `json:"response"`
	Usage    fantasy.Usage `json:"usage"`
	Error    string        `json:"error,omitempty"`
	Sources  []document.SearchResult `json:"sources,omitempty"`
}

type ChatHandler struct {
	providerService *provider.ProviderService
	documentService *document.DocumentService
}

func NewChatHandler(ps *provider.ProviderService, ds *document.DocumentService) *ChatHandler {
	return &ChatHandler{
		providerService: ps,
		documentService: ds,
	}
}

func (h *ChatHandler) HandleChat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	prompt := req.Prompt
	var sources []document.SearchResult

	if req.UseRAG && h.documentService != nil {
		results, err := h.documentService.Search(c.Request.Context(), req.Prompt, 3)
		if err == nil && len(results) > 0 {
			sources = results
			var contextBuilder strings.Builder
			contextBuilder.WriteString("Use the following pieces of context to answer the user's question. If you don't know the answer, just say that you don't know, don't try to make up an answer.\n\n")
			for _, r := range results {
				contextBuilder.WriteString("---\n")
				contextBuilder.WriteString(r.Content)
				contextBuilder.WriteString("\n")
			}
			contextBuilder.WriteString("---\n\nQuestion: ")
			contextBuilder.WriteString(req.Prompt)
			prompt = contextBuilder.String()
		}
	}

	res, err := h.providerService.Chat(c.Request.Context(), req.Provider, req.Model, prompt)
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
		Sources:  sources,
	})
}
