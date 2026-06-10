--! get_submission_by_id
SELECT id, user_id, language, source_code, status, created_at, updated_at
FROM submissions
WHERE id = :id;