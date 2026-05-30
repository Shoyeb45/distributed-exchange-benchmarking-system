-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUserByGithubId :one
SELECT * FROM users
WHERE github_id = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
    name,
    email,
    avatar_url,
    github_id,
    github_username
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;