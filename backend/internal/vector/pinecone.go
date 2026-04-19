package vector

import (
	"context"
	"fmt"

	"github.com/pinecone-io/go-pinecone/pinecone"
	"google.golang.org/protobuf/types/known/structpb"
)

type PineconeStore struct {
	client    *pinecone.Client
	indexName string
	namespace string
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

	return idx, nil
}

func (s *PineconeStore) Add(ctx context.Context, vectors []Vector) error {
	idx, err := s.getIndexConn(ctx)
	if err != nil {
		return err
	}
	defer idx.Close()

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
	idx, err := s.getIndexConn(ctx)
	if err != nil {
		return nil, err
	}
	defer idx.Close()

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
		metadata := match.Vector.Metadata.AsMap()
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
	defer idx.Close()

	return idx.DeleteVectorsById(ctx, ids)
}
