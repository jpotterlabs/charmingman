package vector

import (
	"context"
	"math"
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
		s.vectors[v.ID] = v
	}
	return nil
}

func (s *LocalStore) Search(ctx context.Context, query []float32, limit int) ([]SearchResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	results := make([]SearchResult, 0, len(s.vectors))
	for _, v := range s.vectors {
		score := cosineSimilarity(query, v.Values)
		results = append(results, SearchResult{
			Vector: v,
			Score:  score,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	if len(results) > limit {
		results = results[:limit]
	}

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
