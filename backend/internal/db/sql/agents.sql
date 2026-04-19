-- name: GetAgent :one
SELECT id, name, model, provider, persona, created_at, updated_at
FROM agents
WHERE id = ? LIMIT 1;

-- name: ListAgents :many
SELECT id, name, model, provider, persona, created_at, updated_at
FROM agents
ORDER BY name;

-- name: CreateAgent :one
INSERT INTO agents (
    id, name, model, provider, persona, api_key_ref
) VALUES (
    ?, ?, ?, ?, ?, ?
)
RETURNING id, name, model, provider, persona, created_at, updated_at;

-- name: UpdateAgent :one
UPDATE agents
SET name = ?,
    model = ?,
    provider = ?,
    persona = ?,
    api_key_ref = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING id, name, model, provider, persona, created_at, updated_at;

-- name: DeleteAgent :exec
DELETE FROM agents
WHERE id = ?;

-- name: GetAgentSecret :one
SELECT api_key_ref FROM agents
WHERE id = ? LIMIT 1;
