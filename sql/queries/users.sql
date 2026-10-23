-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password, is_chirpy_red)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3
)
RETURNING *;

-- name: DeleteAllUsers :exec
TRUNCATE TABLE users CASCADE;

-- name: GetUserByEmail :one
SELECT 
    id, 
    created_at, 
    updated_at, 
    email, 
    hashed_password,
    is_chirpy_red
FROM users
WHERE email = $1;

-- name: UpdateUser :exec
UPDATE users
SET
    updated_at = NOW(),
    email = $2,
    hashed_password = $3
WHERE id = $1;

-- name: SetUserChirpyRed :exec
UPDATE users
SET
    updated_at = NOW(),
    is_chirpy_red = $2
WHERE id = $1;