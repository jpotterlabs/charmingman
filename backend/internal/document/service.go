package document

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"charmingman/backend/internal/db"
	"charmingman/backend/internal/vector"
	"github.com/google/uuid"
)

type DocumentService struct {
	queries       *db.Queries
	embedder      vector.Embedder
	vectorStore   vector.VectorStore
	documentsRoot string
}

func NewDocumentService(queries *db.Queries, embedder vector.Embedder, vectorStore vector.VectorStore, documentsRoot string) *DocumentService {
	return &DocumentService{
		queries:       queries,
		embedder:      embedder,
		vectorStore:   vectorStore,
		documentsRoot: documentsRoot,
	}
}

func (s *DocumentService) AddDocument(ctx context.Context, title string, path string) (string, error) {
	// 1. Sanitize and validate path
	cleanPath := filepath.Clean(path)
	
	absRoot, err := filepath.Abs(s.documentsRoot)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute root: %w", err)
	}

	fullPath := filepath.Join(absRoot, cleanPath)
	absPath, err := filepath.Abs(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	rel, err := filepath.Rel(absRoot, absPath)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path escapes documents root: %s", path)
	}

	// 2. Enforce per-file extraction limit (e.g. 10MB)
	fileInfo, err := os.Stat(absPath)
	if err != nil {
		return "", fmt.Errorf("failed to stat file: %w", err)
	}
	if fileInfo.Size() > 10*1024*1024 {
		return "", fmt.Errorf("file too large for extraction (limit 10MB)")
	}

	// 3. Extract text
	text, err := ExtractText(absPath)
	if err != nil {
		return "", fmt.Errorf("failed to extract text: %w", err)
	}

	// 4. Create initial DB record
	docID := uuid.New().String()
	_, err = s.queries.CreateDocument(ctx, db.CreateDocumentParams{
		ID:       docID,
		Title:    title,
		Filename: filepath.Base(cleanPath),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create document record: %w", err)
	}

	// Collected IDs for compensation
	var chunkIDs []string

	// Compensation logic: delete from DB and Vector Store on failure
	cleanup := func() {
		log.Printf("Cleaning up failed document ingestion for docID: %s", docID)
		_ = s.queries.DeleteDocument(context.Background(), docID)
		if len(chunkIDs) > 0 {
			err := s.vectorStore.Delete(context.Background(), chunkIDs)
			if err != nil {
				log.Printf("Warning: failed to delete orphaned vectors for docID %s: %v", docID, err)
			}
		}
	}

	// 5. Chunk and Embed
	chunks := ChunkText(text, 1000, 200)
	chunkTexts := make([]string, len(chunks))
	for i, c := range chunks {
		chunkTexts[i] = c.Content
	}

	embeddings, err := s.embedder.EmbedBatch(ctx, chunkTexts)
	if err != nil {
		cleanup()
		return "", fmt.Errorf("failed to generate embeddings: %w", err)
	}

	if len(embeddings) != len(chunks) {
		cleanup()
		return "", fmt.Errorf("embedding count mismatch: expected %d, got %d", len(chunks), len(embeddings))
	}

	// 6. Save chunks and vectors
	vectors := make([]vector.Vector, len(chunks))
	for i, c := range chunks {
		chunkID := uuid.New().String()
		chunkIDs = append(chunkIDs, chunkID)
		_, err = s.queries.CreateDocumentChunk(ctx, db.CreateDocumentChunkParams{
			ID:          chunkID,
			DocumentID:  docID,
			Content:     c.Content,
			ChunkIndex:  int64(c.Index),
		})
		if err != nil {
			cleanup()
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
		cleanup()
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

	searchParams := make([]SearchResult, 0, len(results))
	for i, r := range results {
		if r.Metadata == nil {
			log.Printf("Warning: nil metadata for match %d (score: %f)", i, r.Score)
			continue
		}
		content, ok1 := r.Metadata["content"].(string)
		docID, ok2 := r.Metadata["document_id"].(string)
		if !ok1 || !ok2 {
			log.Printf("Warning: malformed metadata for match %d (score: %f)", i, r.Score)
			continue
		}
		searchParams = append(searchParams, SearchResult{
			Content:    content,
			DocumentID: docID,
			Score:      r.Score,
		} )
	}

	return searchParams, nil
}
