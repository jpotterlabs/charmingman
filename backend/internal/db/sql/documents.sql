-- name: CreateDocument :one
INSERT INTO documents (
    id, title, filename
) VALUES (
    ?, ?, ?
)
RETURNING *;

-- name: ListDocuments :many
SELECT * FROM documents
ORDER BY created_at DESC;

-- name: GetDocument :one
SELECT * FROM documents
WHERE id = ? LIMIT 1;

-- name: DeleteDocument :exec
DELETE FROM documents
WHERE id = ?;

-- name: CreateDocumentChunk :one
INSERT INTO document_chunks (
    id, document_id, content, chunk_index
) VALUES (
    ?, ?, ?, ?
)
RETURNING *;

-- name: ListChunksByDocument :many
SELECT * FROM document_chunks
WHERE document_id = ?
ORDER BY chunk_index ASC;

-- name: GetChunk :one
SELECT * FROM document_chunks
WHERE id = ? LIMIT 1;
