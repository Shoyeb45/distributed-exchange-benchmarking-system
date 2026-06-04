--! get_submission_by_id : Submission
SELECT id, user_id, problem_id, language, source_code, status,
       runtime_ms, memory_kb, error_message, created_at, updated_at
FROM submissions
WHERE id = :id;