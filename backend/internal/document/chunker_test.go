package document

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChunkText_Basic(t *testing.T) {
	text := "12345678901234567890"
	chunks := ChunkText(text, 10, 2)

	assert.Len(t, chunks, 3)
	assert.Equal(t, "1234567890", chunks[0].Content)
	assert.Equal(t, 0, chunks[0].Index)
	assert.Equal(t, "9012345678", chunks[1].Content)
	assert.Equal(t, 1, chunks[1].Index)
	assert.Equal(t, "7890", chunks[2].Content)
	assert.Equal(t, 2, chunks[2].Index)
}

func TestChunkText_SmallerThanMaxSize(t *testing.T) {
	text := "short text"
	chunks := ChunkText(text, 50, 5)

	assert.Len(t, chunks, 1)
	assert.Equal(t, "short text", chunks[0].Content)
	assert.Equal(t, 0, chunks[0].Index)
}

func TestChunkText_SpaceBreaking(t *testing.T) {
	text := "This is a test of the emergency broadcast system."
	
	// maxSize=20, overlap=5
	// First chunk: "This is a test of th" -> space at 17 -> "This is a test of "
	chunks := ChunkText(text, 20, 5)
	
	assert.NotEmpty(t, chunks)
	assert.Equal(t, "This is a test of", chunks[0].Content)
	
	// Ensure that subsequent chunks are handled correctly
	assert.Contains(t, chunks[1].Content, "the emergency")
}

func TestChunkText_EdgeCases(t *testing.T) {
	// Zero maxSize (should default to len(text))
	chunks := ChunkText("hello", 0, 0)
	assert.Len(t, chunks, 1)
	assert.Equal(t, "hello", chunks[0].Content)

	// Negative overlap (should default to 0)
	chunks = ChunkText("hello", 2, -5)
	assert.Len(t, chunks, 3)
	assert.Equal(t, "he", chunks[0].Content)
	assert.Equal(t, "ll", chunks[1].Content)
	assert.Equal(t, "o", chunks[2].Content)

	// Overlap >= maxSize (should default to maxSize - 1)
	chunks = ChunkText("hello", 3, 5)
	assert.Len(t, chunks, 3)
	assert.Equal(t, "hel", chunks[0].Content)
	assert.Equal(t, "ell", chunks[1].Content)
	assert.Equal(t, "llo", chunks[2].Content)
}
