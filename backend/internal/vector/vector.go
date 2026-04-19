package vector

import (
	"context"
)

// Vector represents a vector in the store.
type Vector struct {
	ID       string
	Values   []float32
	Metadata map[string]interface{}
}

// SearchResult represents a result from a vector search.
type SearchResult struct {
	Vector
	Score float32
}

// VectorStore is the interface for vector databases.
type VectorStore interface {
	Add(ctx context.Context, vectors []Vector) error
	Search(ctx context.Context, query []float32, limit int) ([]SearchResult, error)
	Delete(ctx context.Context, ids []string) error
}

// Embedder is the interface for generating embeddings.
type Embedder interface {
	Embed(ctx context.Context, text string) ([]float32, error)
	EmbedBatch(ctx context.Context, texts []string) ([][]float32, error)
}
