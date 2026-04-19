package vector

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalStore_AddAndSearch(t *testing.T) {
	store := NewLocalStore()
	ctx := context.Background()

	vecs := []Vector{
		{ID: "1", Values: []float32{1.0, 0.0, 0.0}, Metadata: map[string]any{"text": "A"}},
		{ID: "2", Values: []float32{0.0, 1.0, 0.0}, Metadata: map[string]any{"text": "B"}},
		{ID: "3", Values: []float32{0.7, 0.7, 0.0}, Metadata: map[string]any{"text": "C"}},
	}

	err := store.Add(ctx, vecs)
	assert.NoError(t, err)

	// Search for [1.0, 0.0, 0.0]
	results, err := store.Search(ctx, []float32{1.0, 0.0, 0.0}, 2)
	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, "1", results[0].Vector.ID) // Exact match
	assert.Equal(t, "3", results[1].Vector.ID) // Partial match
	assert.InDelta(t, 1.0, results[0].Score, 0.01)
	assert.True(t, results[1].Score > 0 && results[1].Score < 1.0)
}

func TestLocalStore_DeepCopy(t *testing.T) {
	store := NewLocalStore()
	ctx := context.Background()

	originalValues := []float32{1.0, 0.0}
	originalMeta := map[string]any{"key": "value"}
	vec := Vector{ID: "1", Values: originalValues, Metadata: originalMeta}

	err := store.Add(ctx, []Vector{vec})
	assert.NoError(t, err)

	// Mutate original slice and map
	originalValues[0] = 99.0
	originalMeta["key"] = "mutated"

	// Search and verify it wasn't mutated in the store
	results, err := store.Search(ctx, []float32{1.0, 0.0}, 1)
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, float32(1.0), results[0].Vector.Values[0], "vector values should be deep copied")
	assert.Equal(t, "value", results[0].Vector.Metadata["key"], "metadata should be deep copied")
}

func TestLocalStore_Delete(t *testing.T) {
	store := NewLocalStore()
	ctx := context.Background()

	vecs := []Vector{
		{ID: "1", Values: []float32{1.0}},
		{ID: "2", Values: []float32{0.5}},
	}

	err := store.Add(ctx, vecs)
	assert.NoError(t, err)

	err = store.Delete(ctx, []string{"1"})
	assert.NoError(t, err)

	results, err := store.Search(ctx, []float32{1.0}, 2)
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "2", results[0].Vector.ID)
}

func TestLocalStore_CosineSimilarity(t *testing.T) {
	// Identical vectors = 1.0
	score := cosineSimilarity([]float32{1, 0}, []float32{1, 0})
	assert.InDelta(t, 1.0, score, 0.001)

	// Orthogonal vectors = 0.0
	score = cosineSimilarity([]float32{1, 0}, []float32{0, 1})
	assert.InDelta(t, 0.0, score, 0.001)

	// Opposite vectors = -1.0
	score = cosineSimilarity([]float32{1, 0}, []float32{-1, 0})
	assert.InDelta(t, -1.0, score, 0.001)

	// Different lengths = 0.0
	score = cosineSimilarity([]float32{1, 0}, []float32{1, 0, 0})
	assert.Equal(t, float32(0.0), score)
}
