-- name: GetUser :one
SELECT id, name, email, created_at, updated_at 
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, name, email, password_hash, created_at, updated_at 
FROM users
WHERE email = $1;

-- name: ListUsers :many
SELECT id, name, email, created_at, updated_at 
FROM users
ORDER BY id;

-- name: ListUsersByEmail :many
SELECT id, name, email, created_at, updated_at 
FROM users
WHERE email ILIKE $1
ORDER BY id;

-- name: UpsertUser :one
INSERT INTO users (name, email, password_hash, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (email) 
DO UPDATE SET 
    name = EXCLUDED.name,
    password_hash = EXCLUDED.password_hash,
    updated_at = EXCLUDED.updated_at
RETURNING id, name, email, created_at, updated_at;
