package main

import (
	"context"
	"database/sql"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"charmingman/backend/internal/db"
	"charmingman/backend/internal/document"
	"charmingman/backend/internal/handler"
	"charmingman/backend/internal/mcp"
	"charmingman/backend/internal/provider"
	"charmingman/backend/internal/vector"
	"github.com/charmbracelet/openai-go"
	"github.com/charmbracelet/openai-go/option"
	"github.com/charmbracelet/openai-go/packages/param"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type AgentResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Model     string    `json:"model"`
	Provider  string    `json:"provider"`
	Persona   string    `json:"persona"`
	UseRAG    bool      `json:"use_rag"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func listRowToAgentResponse(a db.ListAgentsRow) AgentResponse {
	return AgentResponse{
		ID:        a.ID,
		Name:      a.Name,
		Model:     a.Model,
		Provider:  a.Provider,
		Persona:   a.Persona.String,
		UseRAG:    a.UseRag,
		CreatedAt: a.CreatedAt.Time,
		UpdatedAt: a.UpdatedAt.Time,
	}
}

func createRowToAgentResponse(a db.CreateAgentRow) AgentResponse {
	return AgentResponse{
		ID:        a.ID,
		Name:      a.Name,
		Model:     a.Model,
		Provider:  a.Provider,
		Persona:   a.Persona.String,
		UseRAG:    a.UseRag,
		CreatedAt: a.CreatedAt.Time,
		UpdatedAt: a.UpdatedAt.Time,
	}
}

func main() {
	// Load .env file if it exists
	_ = godotenv.Load()

	ctx := context.Background()

	// Initialize Database
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "."
	}
	dbConn, err := db.Connect(ctx, dataDir)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	queries := db.New(dbConn)

	// Initialize Vector Store and Embedder
	var embedder vector.Embedder
	openaiKey := os.Getenv("OPENAI_API_KEY")
	if openaiKey != "" {
		embedder = vector.NewOpenAIEmbedder(openaiKey, "")
	} else {
		log.Println("Warning: OPENAI_API_KEY not configured. RAG and embeddings will be disabled.")
	}

	var vStore vector.VectorStore = vector.NewLocalStore()
	pineconeKey := os.Getenv("PINECONE_API_KEY")
	pineconeIndex := os.Getenv("PINECONE_INDEX")
	if pineconeKey != "" && pineconeIndex != "" {
		ps, err := vector.NewPineconeStore(pineconeKey, pineconeIndex, "charmingman")
		if err != nil {
			log.Printf("Failed to initialize Pinecone: %v. Falling back to local store.", err)
		} else {
			vStore = ps
			log.Println("Pinecone vector store initialized")
		}
	} else if pineconeKey == "" {
		log.Println("Note: PINECONE_API_KEY not configured. Using local in-memory vector store.")
	}

	// Initialize Services
	docRoot := os.Getenv("DOCUMENTS_ROOT")
	if docRoot == "" {
		docRoot = "./documents"
	}
	if err := os.MkdirAll(docRoot, 0755); err != nil {
		log.Fatalf("Failed to create documents root: %v", err)
	}

	ps := provider.NewProviderService(queries)
	docService := document.NewDocumentService(queries, embedder, vStore, docRoot)

	// Initialize MCP Servers
	mcpServers := os.Getenv("MCP_SERVERS")
	if mcpServers != "" {
		for _, serverCmd := range strings.Split(mcpServers, ",") {
			parts := strings.Fields(serverCmd)
			if len(parts) == 0 {
				continue
			}
			client, err := mcp.NewClient(parts[0], parts[1:]...)
			if err != nil {
				log.Printf("Failed to start MCP server %s: %v", parts[0], err)
				continue
			}
			tools, err := client.ListTools(ctx)
			if err != nil {
				log.Printf("Failed to list tools for MCP server %s: %v", parts[0], err)
				continue
			}
			for _, t := range tools {
				ps.RegisterTool(mcp.NewMCPToolWrapper(client, t))
				log.Printf("Registered MCP tool: %s", t.Name)
			}
		}
	}

	// Register providers from environment variables
	if openaiKey != "" {
		if err := ps.RegisterOpenAI(openaiKey); err != nil {
			log.Printf("Failed to register OpenAI: %v", err)
		} else {
			log.Println("OpenAI provider registered")
		}
	}

	anthropicKey := os.Getenv("ANTHROPIC_API_KEY")
	if anthropicKey != "" {
		if err := ps.RegisterAnthropic(anthropicKey); err != nil {
			log.Printf("Failed to register Anthropic: %v", err)
		} else {
			log.Println("Anthropic provider registered")
		}
	}

	ollamaURL := os.Getenv("OLLAMA_BASE_URL")
	if ollamaURL != "" {
		if err := ps.RegisterLocal("ollama", ollamaURL); err != nil {
			log.Printf("Failed to register Ollama: %v", err)
		} else {
			log.Println("Ollama provider registered at", ollamaURL)
		}
	}

	llamacppURL := os.Getenv("LLAMACPP_BASE_URL")
	if llamacppURL != "" {
		if err := ps.RegisterLocal("llamacpp", llamacppURL); err != nil {
			log.Printf("Failed to register llama.cpp: %v", err)
		} else {
			log.Println("llama.cpp provider registered at", llamacppURL)
		}
	}

	// Initialize Handlers
	chatHandler := handler.NewChatHandler(ps, docService, queries)
	
	var transcribeHandler *handler.TranscribeHandler
	if openaiKey != "" {
		transcribeHandler = handler.NewTranscribeHandler(openaiKey, "")
	}

	// Setup Gin router
	r := gin.Default()

	// Unified API Key Middleware
	r.Use(func(c *gin.Context) {
		gatewayKey := os.Getenv("GATEWAY_API_KEY")
		if gatewayKey == "" {
			c.Next()
			return
		}

		apiKey := c.GetHeader("X-Charming-Key")
		if apiKey != gatewayKey {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	})

	// Routes
	v1 := r.Group("/api/v1")
	{
		v1.POST("/chat", chatHandler.HandleChat)
		
		v1.GET("/usage", func(c *gin.Context) {
			limit, err := strconv.ParseInt(c.DefaultQuery("limit", "100"), 10, 64)
			if err != nil || limit <= 0 || limit > 1000 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit. Must be a positive integer <= 1000"})
				return
			}

			offset, err := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)
			if err != nil || offset < 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset. Must be a non-negative integer"})
				return
			}
			
			logs, err := queries.ListUsageLogs(c.Request.Context(), db.ListUsageLogsParams{
				Limit:  limit,
				Offset: offset,
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			
			stats, err := queries.GetTotalUsage(c.Request.Context())
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			
			c.JSON(http.StatusOK, gin.H{
				"logs":  logs,
				"stats": stats,
			})
		})

		v1.GET("/agents", func(c *gin.Context) {
			agents, err := queries.ListAgents(c.Request.Context())
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			
			resp := make([]AgentResponse, len(agents))
			for i, a := range agents {
				resp[i] = listRowToAgentResponse(a)
			}
			c.JSON(http.StatusOK, resp)
		})

		v1.POST("/agents", func(c *gin.Context) {
			var req struct {
				ID       string `json:"id"`
				Name     string `json:"name" binding:"required"`
				Model    string `json:"model" binding:"required"`
				Provider string `json:"provider" binding:"required"`
				Persona  string `json:"persona"`
				APIKey   string `json:"api_key"`
				UseRAG   bool   `json:"use_rag"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			if req.ID == "" {
				req.ID = uuid.New().String()
			}

			agent, err := queries.CreateAgent(c.Request.Context(), db.CreateAgentParams{
				ID:        req.ID,
				Name:      req.Name,
				Model:     req.Model,
				Provider:  req.Provider,
				Persona:   sql.NullString{String: req.Persona, Valid: req.Persona != ""},
				ApiKeyRef: sql.NullString{String: req.APIKey, Valid: req.APIKey != ""},
				UseRag:    req.UseRAG,
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusCreated, createRowToAgentResponse(agent))
		})

		// Document Endpoints
		v1.POST("/documents", func(c *gin.Context) {
			c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1024*10) 

			var req struct {
				Title string `json:"title" binding:"required"`
				Path  string `json:"path" binding:"required"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			if embedder == nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Embedder not initialized (OPENAI_API_KEY missing)"})
				return
			}

			docID, err := docService.AddDocument(c.Request.Context(), req.Title, req.Path)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusCreated, gin.H{"id": docID})
		})

		v1.GET("/search", func(c *gin.Context) {
			query := c.Query("q")
			if query == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
				return
			}

			limitStr := c.DefaultQuery("limit", "5")
			limit, err := strconv.Atoi(limitStr)
			if err != nil || limit <= 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit. Must be a positive integer"})
				return
			}
			if limit > 100 {
				limit = 100
			}

			if embedder == nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Embedder not initialized"})
				return
			}

			results, err := docService.Search(c.Request.Context(), query, limit)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, results)
		})

		// Multimedia Endpoints
		if transcribeHandler != nil {
			v1.POST("/transcribe", transcribeHandler.HandleTranscribe)
		}

		v1.POST("/speech", func(c *gin.Context) {
			var req struct {
				Text string `json:"text" binding:"required"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			if openaiKey == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "OpenAI API key not configured"})
				return
			}

			client := openai.NewClient(option.WithAPIKey(openaiKey))
			resp, err := client.Audio.Speech.New(ctx, openai.AudioSpeechNewParams{
				Input: req.Text,
				Model: openai.SpeechModelTTS1,
				Voice: openai.AudioSpeechNewParamsVoiceUnion{
					OfString: param.NewOpt("alloy"),
				},
				ResponseFormat: openai.AudioSpeechNewParamsResponseFormatMP3,
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			defer resp.Body.Close()

			c.Header("Content-Type", "audio/mpeg")
			io.Copy(c.Writer, resp.Body)
		})

		// Tooling Endpoints
		v1.GET("/tools", func(c *gin.Context) {
			tools := ps.ListTools()
			var resp []map[string]interface{}
			for _, t := range tools {
				info := t.Info()
				resp = append(resp, map[string]interface{}{
					"name":        info.Name,
					"description": info.Description,
				})
			}
			c.JSON(http.StatusOK, resp)
		})
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}

	log.Printf("AI Gateway starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
