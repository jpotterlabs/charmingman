package agent

import (
	"context"

	"charm.land/fantasy"
	"charm.land/fantasy/providers/openai"
)

type Agent struct {
	fantasyAgent *fantasy.Agent
}

func NewAgent(apiKey string, modelName string) (*Agent, error) {
	provider, err := openai.New(openai.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	model, err := provider.LanguageModel(ctx, modelName)
	if err != nil {
		return nil, err
	}

	fantasyAgent := fantasy.NewAgent(
		model,
		fantasy.WithSystemPrompt("You are a helpful assistant named CharmingMan."),
	)

	return &Agent{
		fantasyAgent: fantasyAgent,
	}, nil
}

func (a *Agent) Generate(ctx context.Context, prompt string) (string, error) {
	result, err := a.fantasyAgent.Generate(ctx, fantasy.AgentCall{Prompt: prompt})
	if err != nil {
		return "", err
	}
	return result.Response.Content.Text(), nil
}
