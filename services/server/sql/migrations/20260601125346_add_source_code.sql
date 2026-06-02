-- +goose Up
ALTER TABLE submissions
ADD COLUMN source_code TEXT DEFAULT '';

-- +goose Down

ALTER TABLE submissions
DROP COLUMN source_code;

