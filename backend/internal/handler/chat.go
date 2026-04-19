package handler

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"strings"

	"charmingman/backend/internal/db"
	"charmingman/backend/internal/document"
	"charmingman/backend/internal/provider"
	"charm.land/fantasy"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChatRequest struct {
	Provider string `json:"provider" binding:"required"`
	Model    string `json:"model" binding:"required"`
	Prompt   string `json:"prompt" binding:"required"`
	UseRAG   bool   `json:"use_rag"`
	RoomID   string `json:"room_id"`
	AgentID  string `json:"agent_id"`
}

type ChatResponse struct {
	Response string                  `json:"response"`
	Usage    fantasy.Usage           `json:"usage"`
	Error    string                  `json:"error,omitempty"`
	Sources  []document.SearchResult `json:"sources,omitempty"`
}

type ChatHandler struct {
	providerService *provider.ProviderService
	documentService *document.DocumentService
	queries         *db.Queries
}

func NewChatHandler(ps *provider.ProviderService, ds *document.DocumentService, q *db.Queries) *ChatHandler {
	return &ChatHandler{
		providerService: ps,
		documentService: ds,
		queries:         q,
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

	// 1. Fetch History if RoomID provided
	var history []fantasy.Message
	if req.RoomID != "" && h.queries != nil {
		dbMsgs, err := h.queries.ListMessagesByRoom(c.Request.Context(), req.RoomID)
		if err == nil {
			for _, m := range dbMsgs {
				role := fantasy.MessageRole(m.Role)
				history = append(history, fantasy.Message{
					Role:    role,
					Content: []fantasy.MessagePart{fantasy.TextPart{Text: m.Content}},
				})
			}
		}
	}

	// 2. Perform RAG Search if requested
	if req.UseRAG && h.documentService != nil {
		results, err := h.documentService.Search(c.Request.Context(), req.Prompt, 3)
		if err != nil {
			safePrompt := req.Prompt
			if len(safePrompt) > 20 {
				safePrompt = safePrompt[:20] + "..."
			}
			log.Printf("Error during RAG search for prompt %q (len: %d): %v", safePrompt, len(req.Prompt), err)
		} else if len(results) > 0 {
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

	// 3. Persist User Message
	if req.RoomID != "" && h.queries != nil {
		_, _ = h.queries.CreateMessage(context.Background(), db.CreateMessageParams{
			ID:      uuid.New().String(),
			RoomID:  req.RoomID,
			Role:    string(fantasy.MessageRoleUser),
			Content: req.Prompt,
		})
	}

	// 4. Call Model
	res, err := h.providerService.Chat(c.Request.Context(), req.Provider, req.Model, prompt, history)
	if err != nil {
		if strings.Contains(err.Error(), "not registered") {
			c.JSON(http.StatusBadRequest, ChatResponse{Error: err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, ChatResponse{Error: err.Error()})
		}
		return
	}

	// 5. Persist Assistant Message
	if req.RoomID != "" && h.queries != nil {
		agentID := sql.NullString{String: req.AgentID, Valid: req.AgentID != ""}
		_, _ = h.queries.CreateMessage(context.Background(), db.CreateMessageParams{
			ID:         uuid.New().String(),
			RoomID:     req.RoomID,
			AgentID:    agentID,
			Role:       string(fantasy.MessageRoleAssistant),
			Content:    res.Content.Text(),
			TokensUsed: sql.NullInt64{Int64: res.Usage.TotalTokens, Valid: true},
		})
	}

	c.JSON(http.StatusOK, ChatResponse{
		Response: res.Content.Text(),
		Usage:    res.Usage,
		Sources:  sources,
	})
}
