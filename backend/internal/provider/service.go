package provider

import (
	"context"
	"fmt"
	"time"

	"charm.land/fantasy"
	"charm.land/fantasy/providers/anthropic"
	"charm.land/fantasy/providers/openai"
	"charm.land/fantasy/providers/openaicompat"
)

type ProviderService struct {
	// Map of provider name to fantasy provider
	providers map[string]fantasy.Provider
}

func NewProviderService() *ProviderService {
	return &ProviderService{
		providers: make(map[string]fantasy.Provider),
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

func (s *ProviderService) Chat(ctx context.Context, providerName, modelName, prompt string) (*fantasy.Response, error) {
	p, ok := s.providers[providerName]
	if !ok {
		return nil, fmt.Errorf("provider %s not registered", providerName)
	}

	// Set 30s timeout for model calls
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	model, err := p.LanguageModel(ctx, modelName)
	if err != nil {
		return nil, err
	}

	agent := fantasy.NewAgent(model)
	result, err := agent.Generate(ctx, fantasy.AgentCall{Prompt: prompt})
	if err != nil {
		return nil, err
	}

	return &result.Response, nil
}
