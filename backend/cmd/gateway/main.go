package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"charmingman/backend/internal/db"
	"charmingman/backend/internal/document"
	"charmingman/backend/internal/handler"
	"charmingman/backend/internal/provider"
	"charmingman/backend/internal/vector"
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
	}

	// Initialize Services
	ps := provider.NewProviderService(queries)
	docService := document.NewDocumentService(queries, embedder, vStore)

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
	chatHandler := handler.NewChatHandler(ps, docService)

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
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			if req.ID == "" {
				req.ID = uuid.New().String()
			}

			agent, err := queries.CreateAgent(c.Request.Context(), db.CreateAgentParams{
				ID:       req.ID,
				Name:     req.Name,
				Model:    req.Model,
				Provider: req.Provider,
				Persona:  sql.NullString{String: req.Persona, Valid: req.Persona != ""},
				ApiKeyRef: sql.NullString{String: req.APIKey, Valid: req.APIKey != ""},
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusCreated, createRowToAgentResponse(agent))
		})

		// Document Endpoints
		v1.POST("/documents", func(c *gin.Context) {
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
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
			limit, _ := strconv.Atoi(limitStr)

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
