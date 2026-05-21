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


-- KEYSTORES
CREATE TABLE keystores (
    id SERIAL PRIMARY KEY,
    client_id INT NOT NULL,
    primary_key TEXT NOT NULL,
    secondary_key TEXT NOT NULL,
    status BOOLEAN DEFAULT TRUE,

    refresh_token TEXT DEFAULT '',
    device_fingerprint TEXT DEFAULT '',

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),

    CONSTRAINT fk_keystores_user
        FOREIGN KEY (client_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE TRIGGER set_updated_at_keystores
BEFORE UPDATE ON keystores
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE INDEX idx_keystores_client_id ON keystores(client_id);
CREATE INDEX idx_keystores_client_primary_status ON keystores(client_id, primary_key, status);
CREATE INDEX idx_keystores_client_primary_secondary ON keystores(client_id, primary_key, secondary_key);

-- +goose Down

DROP TABLE IF EXISTS keystores;
DROP TABLE IF EXISTS users;