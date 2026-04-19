package vector

import (
	"context"
	"fmt"
	"sync"

	"github.com/pinecone-io/go-pinecone/pinecone"
	"google.golang.org/protobuf/types/known/structpb"
)

type PineconeStore struct {
	client    *pinecone.Client
	indexName string
	namespace string

	// Connection cache
	mu           sync.RWMutex
	resolvedHost string
	cachedIdx    *pinecone.IndexConnection
}

func NewPineconeStore(apiKey string, indexName string, namespace string) (*PineconeStore, error) {
	pc, err := pinecone.NewClient(pinecone.NewClientParams{
		ApiKey: apiKey,
	})
	if err != nil {
		return nil, err
	}

	return &PineconeStore{
		client:    pc,
		indexName: indexName,
		namespace: namespace,
	}, nil
}

func (s *PineconeStore) getIndexConn(ctx context.Context) (*pinecone.IndexConnection, error) {
	s.mu.RLock()
	if s.cachedIdx != nil {
		defer s.mu.RUnlock()
		return s.cachedIdx, nil
	}
	s.mu.RUnlock()

	s.mu.Lock()
	defer s.mu.Unlock()

	// Double check after acquiring write lock
	if s.cachedIdx != nil {
		return s.cachedIdx, nil
	}

	idxMeta, err := s.client.DescribeIndex(ctx, s.indexName)
	if err != nil {
		return nil, fmt.Errorf("failed to describe index: %w", err)
	}

	idx, err := s.client.Index(pinecone.NewIndexConnParams{
		Host:      idxMeta.Host,
		Namespace: s.namespace,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to index: %w", err)
	}

	s.resolvedHost = idxMeta.Host
	s.cachedIdx = idx
	return idx, nil
}

func (s *PineconeStore) Add(ctx context.Context, vectors []Vector) error {
	idx, err := s.getIndexConn(ctx)
	if err != nil {
		return err
	}

	pineconeVectors := make([]*pinecone.Vector, len(vectors))
	for i, v := range vectors {
		metadata, err := structpb.NewStruct(v.Metadata)
		if err != nil {
			return fmt.Errorf("failed to convert metadata: %w", err)
		}

		pineconeVectors[i] = &pinecone.Vector{
			Id:       v.ID,
			Values:   v.Values,
			Metadata: metadata,
		}
	}

	_, err = idx.UpsertVectors(ctx, pineconeVectors)
	return err
}

func (s *PineconeStore) Search(ctx context.Context, query []float32, limit int) ([]SearchResult, error) {
	if limit <= 0 {
		return nil, nil
	}
	if limit > 10000 {
		limit = 10000
	}

	idx, err := s.getIndexConn(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := idx.QueryByVectorValues(ctx, &pinecone.QueryByVectorValuesRequest{
		Vector:          query,
		TopK:            uint32(limit),
		IncludeMetadata: true,
	})
	if err != nil {
		return nil, err
	}

	results := make([]SearchResult, len(resp.Matches))
	for i, match := range resp.Matches {
		var metadata map[string]interface{}
		if match.Vector != nil && match.Vector.Metadata != nil {
			metadata = match.Vector.Metadata.AsMap()
		} else {
			metadata = make(map[string]interface{})
		}
		
		results[i] = SearchResult{
			Vector: Vector{
				ID:       match.Vector.Id,
				Values:   match.Vector.Values,
				Metadata: metadata,
			},
			Score: match.Score,
		}
	}

	return results, nil
}

func (s *PineconeStore) Delete(ctx context.Context, ids []string) error {
	idx, err := s.getIndexConn(ctx)
	if err != nil {
		return err
	}

	return idx.DeleteVectorsById(ctx, ids)
}

// Close closes the cached index connection if it exists.
func (s *PineconeStore) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.cachedIdx != nil {
		err := s.cachedIdx.Close()
		s.cachedIdx = nil
		return err
	}
	return nil
}
