-- name: ListUsers :many
SELECT id, name, email, created_at, updated_at 
FROM users
ORDER BY id;

-- name: UpsertUser :one
INSERT INTO users (name, email, created_at, updated_at)
VALUES ($1, $2, $3, $4)
ON CONFLICT (email) 
DO UPDATE SET 
    name = EXCLUDED.name,
    updated_at = EXCLUDED.updated_at
RETURNING id, name, email, created_at, updated_at;
