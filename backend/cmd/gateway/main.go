package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"

	"charmingman/backend/internal/db"
	"charmingman/backend/internal/handler"
	"charmingman/backend/internal/provider"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

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

	// Initialize Provider Service
	ps := provider.NewProviderService(queries)

	// Register providers from environment variables
	openaiKey := os.Getenv("OPENAI_API_KEY")
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
	chatHandler := handler.NewChatHandler(ps)

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
			limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "100"), 10, 64)
			offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)
			
			logs, err := queries.ListUsageLogs(c.Request.Context(), db.ListUsageLogsParams{
				Limit:  limit,
				Offset: offset,
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			
			stats, _ := queries.GetTotalUsage(c.Request.Context())
			
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
			c.JSON(http.StatusOK, agents)
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
				ApiKey:   sql.NullString{String: req.APIKey, Valid: req.APIKey != ""},
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusCreated, agent)
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
