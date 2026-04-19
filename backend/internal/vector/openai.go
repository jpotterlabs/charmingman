package vector

import (
	"context"
	"fmt"

	"github.com/charmbracelet/openai-go"
	"github.com/charmbracelet/openai-go/option"
	"github.com/charmbracelet/openai-go/packages/param"
)

type OpenAIEmbedder struct {
	client *openai.Client
	model  string
}

func NewOpenAIEmbedder(apiKey string, model string) *OpenAIEmbedder {
	if model == "" {
		model = "text-embedding-3-small"
	}
	c := openai.NewClient(option.WithAPIKey(apiKey))
	return &OpenAIEmbedder{
		client: &c,
		model:  model,
	}
}

func (e *OpenAIEmbedder) Embed(ctx context.Context, text string) ([]float32, error) {
	resp, err := e.client.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Input: openai.EmbeddingNewParamsInputUnion{
			OfString: param.NewOpt(text),
		},
		Model: e.model,
	})
	if err != nil {
		return nil, err
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no embedding returned")
	}

	f64s := resp.Data[0].Embedding
	f32s := make([]float32, len(f64s))
	for i, f := range f64s {
		f32s[i] = float32(f)
	}

	return f32s, nil
}

func (e *OpenAIEmbedder) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
	resp, err := e.client.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Input: openai.EmbeddingNewParamsInputUnion{
			OfArrayOfStrings: texts,
		},
		Model: e.model,
	})
	if err != nil {
		return nil, err
	}

	results := make([][]float32, len(resp.Data))
	for i, data := range resp.Data {
		f32s := make([]float32, len(data.Embedding))
		for j, f := range data.Embedding {
			f32s[j] = float32(f)
		}
		results[i] = f32s
	}

	return results, nil
}
