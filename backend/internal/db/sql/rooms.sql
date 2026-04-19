-- name: GetRoom :one
SELECT * FROM rooms
WHERE id = ? LIMIT 1;

-- name: ListRooms :many
SELECT * FROM rooms
ORDER BY created_at DESC;

-- name: CreateRoom :one
INSERT INTO rooms (
    id, title, goal
) VALUES (
    ?, ?, ?
)
RETURNING *;

-- name: DeleteRoom :exec
DELETE FROM rooms
WHERE id = ?;
