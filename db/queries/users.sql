-- name: CreateUser :one
INSERT INTO users (
    full_name, email, password, role
) VALUES (
             $1, $2, $3, $4
         ) RETURNING id, full_name, email, role, created_at;

-- name: GetUserByEmail :one
SELECT id, full_name, email, password, role
FROM users
WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT id, full_name, email, role
FROM users
ORDER BY id
    LIMIT $1 OFFSET $2;