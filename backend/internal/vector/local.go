package vector

import (
	"context"
	"math"
	"maps"
	"slices"
	"sort"
	"sync"
)

type LocalStore struct {
	vectors map[string]Vector
	mu      sync.RWMutex
}

func NewLocalStore() *LocalStore {
	return &LocalStore{
		vectors: make(map[string]Vector),
	}
}

func (s *LocalStore) Add(ctx context.Context, vectors []Vector) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, v := range vectors {
		// Deep copy to avoid leaking internal state
		copiedValues := slices.Clone(v.Values)
		copiedMetadata := maps.Clone(v.Metadata)
		
		s.vectors[v.ID] = Vector{
			ID:       v.ID,
			Values:   copiedValues,
			Metadata: copiedMetadata,
		}
	}
	return nil
}

func (s *LocalStore) Search(ctx context.Context, query []float32, limit int) ([]SearchResult, error) {
	if limit <= 0 {
		return nil, nil
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	results := make([]SearchResult, 0, len(s.vectors))
	for _, v := range s.vectors {
		score := cosineSimilarity(query, v.Values)
		
		// Deep copy result to avoid leaking internal state
		copiedValues := slices.Clone(v.Values)
		copiedMetadata := maps.Clone(v.Metadata)

		results = append(results, SearchResult{
			Vector: Vector{
				ID:       v.ID,
				Values:   copiedValues,
				Metadata: copiedMetadata,
			},
			Score:  score,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	if limit > len(results) {
		limit = len(results)
	}
	results = results[:limit]

	return results, nil
}

func (s *LocalStore) Delete(ctx context.Context, ids []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, id := range ids {
		delete(s.vectors, id)
	}
	return nil
}

func cosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) {
		return 0
	}
	var dotProduct, normA, normB float32
	for i := 0; i < len(a); i++ {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dotProduct / (float32(math.Sqrt(float64(normA))) * float32(math.Sqrt(float64(normB))))
}
