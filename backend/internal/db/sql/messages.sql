-- name: ListMessagesByRoom :many
SELECT * FROM messages
WHERE room_id = ?
ORDER BY created_at ASC, id ASC;

-- name: CreateMessage :one
INSERT INTO messages (
    id, room_id, agent_id, role, content, tokens_used
) VALUES (
    ?, ?, ?, ?, ?, ?
)
RETURNING *;
