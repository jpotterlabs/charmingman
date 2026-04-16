package main

import (
	"log"
	"net/http"
	"os"

	"charmingman/backend/internal/handler"
	"charmingman/backend/internal/provider"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists
	_ = godotenv.Load()

	// Initialize Provider Service
	ps := provider.NewProviderService()

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
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("AI Gateway starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
