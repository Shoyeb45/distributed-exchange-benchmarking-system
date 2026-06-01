-- name: CreateSubmission :one
INSERT INTO submissions (
    user_id,
    language,
    source_code,
    status
)
VALUES (
    $1,
    $2,
    $3,
    'UPLOADED'
)
RETURNING *;