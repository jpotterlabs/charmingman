package document

import (
	"strings"
)

// Chunk represents a piece of text from a document.
type Chunk struct {
	Content string
	Index   int
}

// ChunkText splits a large string into smaller chunks with overlap.
func ChunkText(text string, maxSize int, overlap int) []Chunk {
	if len(text) <= maxSize {
		return []Chunk{{Content: text, Index: 0}}
	}

	var chunks []Chunk
	index := 0
	
	for start := 0; start < len(text); {
		end := start + maxSize
		if end >= len(text) {
			end = len(text)
		} else {
			// Try to find a good breaking point (newline or space)
			lastNewline := strings.LastIndex(text[start:end], "\n")
			if lastNewline != -1 && lastNewline > maxSize/2 {
				end = start + lastNewline + 1
			} else {
				lastSpace := strings.LastIndex(text[start:end], " ")
				if lastSpace != -1 && lastSpace > maxSize/2 {
					end = start + lastSpace + 1
				}
			}
		}

		chunks = append(chunks, Chunk{
			Content: strings.TrimSpace(text[start:end]),
			Index:   index,
		})
		index++

		if end == len(text) {
			break
		}

		start = end - overlap
		if start < 0 {
			start = 0
		}
		// Guard against infinite loop if maxSize is too small for overlap
		if start <= (end - maxSize) {
			start = end
		}
	}

	return chunks
}
