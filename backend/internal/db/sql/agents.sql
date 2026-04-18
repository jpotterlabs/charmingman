-- name: GetAgent :one
SELECT * FROM agents
WHERE id = ? LIMIT 1;

-- name: ListAgents :many
SELECT * FROM agents
ORDER BY name;

-- name: CreateAgent :one
INSERT INTO agents (
    id, name, model, provider, persona, api_key
) VALUES (
    ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: UpdateAgent :one
UPDATE agents
SET name = ?,
    model = ?,
    provider = ?,
    persona = ?,
    api_key = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: DeleteAgent :exec
DELETE FROM agents
WHERE id = ?;
