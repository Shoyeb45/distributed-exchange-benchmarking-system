-- +goose Up

-- function to automatic update the table if something is getting updated
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd
-- USERS
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    avatar_url TEXT NOT NULL,
    github_id INT NOT NULL,
    github_username TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TRIGGER set_updated_at_users
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_github_username ON users(github_username);


-- refresh_tokens
CREATE TABLE refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TRIGGER set_updated_at_refresh_tokens
BEFORE UPDATE ON refresh_tokens
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- CREATE INDEX idx_refresh_tokens_client_id ON refresh_tokens(client_id);
-- CREATE INDEX idx_refresh_tokens_client_primary_status ON refresh_tokens(client_id, primary_key, status);
-- CREATE INDEX idx_refresh_tokens_client_primary_secondary ON refresh_tokens(client_id, primary_key, secondary_key);

-- +goose Down

DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS users;