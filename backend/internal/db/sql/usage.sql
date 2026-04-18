-- name: LogUsage :one
INSERT INTO usage_log (
    provider, model, prompt_tokens, completion_tokens, total_tokens, latency_ms, cost
) VALUES (
    ?, ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: GetTotalUsage :one
SELECT
    SUM(total_tokens) as total_tokens,
    SUM(cost) as total_cost,
    COUNT(*) as total_requests
FROM usage_log;

-- name: ListUsageLogs :many
SELECT * FROM usage_log
ORDER BY timestamp DESC
LIMIT ? OFFSET ?;
