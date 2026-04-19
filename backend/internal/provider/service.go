package provider

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"charmingman/backend/internal/db"
	"charm.land/fantasy"
	"charm.land/fantasy/providers/anthropic"
	"charm.land/fantasy/providers/openai"
	"charm.land/fantasy/providers/openaicompat"
)

type ProviderService struct {
	// Map of provider name to fantasy provider
	providers map[string]fantasy.Provider
	queries   *db.Queries
}

func NewProviderService(queries *db.Queries) *ProviderService {
	return &ProviderService{
		providers: make(map[string]fantasy.Provider),
		queries:   queries,
	}
}

func (s *ProviderService) RegisterOpenAI(apiKey string) error {
	p, err := openai.New(openai.WithAPIKey(apiKey))
	if err != nil {
		return err
	}
	s.providers["openai"] = p
	return nil
}

func (s *ProviderService) RegisterAnthropic(apiKey string) error {
	p, err := anthropic.New(anthropic.WithAPIKey(apiKey))
	if err != nil {
		return err
	}
	s.providers["anthropic"] = p
	return nil
}

func (s *ProviderService) RegisterLocal(name, baseURL string) error {
	p, err := openaicompat.New(
		openaicompat.WithName(name),
		openaicompat.WithBaseURL(baseURL),
		openaicompat.WithAPIKey("ollama"), // Ollama doesn't need a real key but some SDKs expect one
	)
	if err != nil {
		return err
	}
	s.providers[name] = p
	return nil
}

func (s *ProviderService) Chat(ctx context.Context, providerName, modelName string, history []fantasy.Message) (*fantasy.Response, error) {
	p, ok := s.providers[providerName]
	if !ok {
		return nil, fmt.Errorf("provider %s not registered", providerName)
	}

	// Set 30s timeout for model calls
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	startTime := time.Now()
	model, err := p.LanguageModel(ctx, modelName)
	if err != nil {
		return nil, err
	}

	agent := fantasy.NewAgent(model)
	
	call := fantasy.AgentCall{
		Messages: history,
	}

	result, err := agent.Generate(ctx, call)
	if err != nil {
		return nil, err
	}
	latency := time.Since(startTime)

	// Log usage asynchronously with a short timeout
	if s.queries != nil {
		go func() {
			logCtx, logCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer logCancel()

			_, err := s.queries.LogUsage(logCtx, db.LogUsageParams{
				Provider:         providerName,
				Model:            modelName,
				PromptTokens:     int64(result.Response.Usage.InputTokens),
				CompletionTokens: int64(result.Response.Usage.OutputTokens),
				TotalTokens:      int64(result.Response.Usage.TotalTokens),
				LatencyMs:        latency.Milliseconds(),
				Cost:             sql.NullFloat64{Valid: false}, // Unknown cost for now
			})
			if err != nil {
				log.Printf("Warning: failed to log usage: %v", err)
			}
		}()
	}

	return &result.Response, nil
}
