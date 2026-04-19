package document

import (
	"context"
	"fmt"

	"charmingman/backend/internal/db"
	"charmingman/backend/internal/vector"
	"github.com/google/uuid"
)

type DocumentService struct {
	queries     *db.Queries
	embedder    vector.Embedder
	vectorStore vector.VectorStore
}

func NewDocumentService(queries *db.Queries, embedder vector.Embedder, vectorStore vector.VectorStore) *DocumentService {
	return &DocumentService{
		queries:     queries,
		embedder:    embedder,
		vectorStore: vectorStore,
	}
}

func (s *DocumentService) AddDocument(ctx context.Context, title string, path string) (string, error) {
	text, err := ExtractText(path)
	if err != nil {
		return "", fmt.Errorf("failed to extract text: %w", err)
	}

	docID := uuid.New().String()
	_, err = s.queries.CreateDocument(ctx, db.CreateDocumentParams{
		ID:       docID,
		Title:    title,
		Filename: path,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create document record: %w", err)
	}

	chunks := ChunkText(text, 1000, 200)
	chunkTexts := make([]string, len(chunks))
	for i, c := range chunks {
		chunkTexts[i] = c.Content
	}

	embeddings, err := s.embedder.EmbedBatch(ctx, chunkTexts)
	if err != nil {
		return "", fmt.Errorf("failed to generate embeddings: %w", err)
	}

	vectors := make([]vector.Vector, len(chunks))
	for i, c := range chunks {
		chunkID := uuid.New().String()
		_, err = s.queries.CreateDocumentChunk(ctx, db.CreateDocumentChunkParams{
			ID:          chunkID,
			DocumentID:  docID,
			Content:     c.Content,
			ChunkIndex:  int64(c.Index),
		})
		if err != nil {
			return "", fmt.Errorf("failed to save chunk %d: %w", i, err)
		}

		vectors[i] = vector.Vector{
			ID:     chunkID,
			Values: embeddings[i],
			Metadata: map[string]interface{}{
				"document_id": docID,
				"content":     c.Content,
			},
		}
	}

	err = s.vectorStore.Add(ctx, vectors)
	if err != nil {
		return "", fmt.Errorf("failed to add vectors to store: %w", err)
	}

	return docID, nil
}

type SearchResult struct {
	Content    string  `json:"content"`
	DocumentID string  `json:"document_id"`
	Score      float32 `json:"score"`
}

func (s *DocumentService) Search(ctx context.Context, query string, limit int) ([]SearchResult, error) {
	emb, err := s.embedder.Embed(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to embed query: %w", err)
	}

	results, err := s.vectorStore.Search(ctx, emb, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search vector store: %w", err)
	}

	searchParams := make([]SearchResult, len(results))
	for i, r := range results {
		content, _ := r.Metadata["content"].(string)
		docID, _ := r.Metadata["document_id"].(string)
		searchParams[i] = SearchResult{
			Content:    content,
			DocumentID: docID,
			Score:      r.Score,
		}
	}

	return searchParams, nil
}
