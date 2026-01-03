CREATE TABLE IF NOT EXISTS oauth_accounts (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) NOT NULL,
    provider TEXT NOT NULL,
    provider_user_id TEXT NOT NULL,
    changed_date TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP, 
    create_date TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(provider, provider_user_id)
);

ALTER TABLE users ALTER COLUMN password_hash DROP NOT NULL;