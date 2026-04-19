-- name: LogUsage :one
INSERT INTO usage_log (
    provider, model, prompt_tokens, completion_tokens, total_tokens, latency_ms, cost
) VALUES (
    ?, ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: GetTotalUsage :one
SELECT
    CAST(COALESCE(SUM(total_tokens), 0) AS INTEGER) as total_tokens,
    CAST(COALESCE(SUM(cost), 0.0) AS REAL) as total_cost,
    COUNT(*) as total_requests
FROM usage_log;

-- name: ListUsageLogs :many
SELECT * FROM usage_log
ORDER BY timestamp DESC, id DESC
LIMIT ? OFFSET ?;
