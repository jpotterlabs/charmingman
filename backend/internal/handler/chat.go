package handler

import (
	"database/sql"
	"log"
	"net/http"
	"slices"
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

	// 1. Fetch History if RoomID provided (Bounded to last 10 messages)
	var history []fantasy.Message
	if req.RoomID != "" && h.queries != nil {
		dbMsgs, err := h.queries.ListMessagesByRoom(c.Request.Context(), req.RoomID)
		if err == nil {
			// Implementation of fixed cap (last 10 messages)
			startIdx := 0
			if len(dbMsgs) > 10 {
				startIdx = len(dbMsgs) - 10
			}
			for i := startIdx; i < len(dbMsgs); i++ {
				m := dbMsgs[i]
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
			// Redact prompt in logs completely, only log metadata
			log.Printf("Error during RAG search for prompt of len %d: %v", len(req.Prompt), err)
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

	// 3. Persist User Message (Respect Request Context)
	if req.RoomID != "" && h.queries != nil {
		_, _ = h.queries.CreateMessage(c.Request.Context(), db.CreateMessageParams{
			ID:      uuid.New().String(),
			RoomID:  req.RoomID,
			Role:    string(fantasy.MessageRoleUser),
			Content: req.Prompt,
		})
	}

	// 4. Augment history with current prompt before calling model
	augmentedHistory := append(slices.Clone(history), fantasy.Message{
		Role:    fantasy.MessageRoleUser,
		Content: []fantasy.MessagePart{fantasy.TextPart{Text: prompt}},
	})

	// 5. Call Model (using only history which now contains the current prompt)
	res, err := h.providerService.Chat(c.Request.Context(), req.Provider, req.Model, augmentedHistory)
	if err != nil {
		if strings.Contains(err.Error(), "not registered") {
			c.JSON(http.StatusBadRequest, ChatResponse{Error: err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, ChatResponse{Error: err.Error()})
		}
		return
	}

	// 6. Persist Assistant Message (Respect Request Context)
	if req.RoomID != "" && h.queries != nil {
		agentID := sql.NullString{String: req.AgentID, Valid: req.AgentID != ""}
		_, _ = h.queries.CreateMessage(c.Request.Context(), db.CreateMessageParams{
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
