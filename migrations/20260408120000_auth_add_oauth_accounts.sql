-- +goose Up
-- +goose StatementBegin
CREATE TABLE auth.oauth_accounts (
    id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR NOT NULL,
    provider VARCHAR NOT NULL,
    provider_user_id VARCHAR NOT NULL,
    provider_email VARCHAR,
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX uk_oauth_accounts_provider_user
    ON auth.oauth_accounts (provider, provider_user_id);

CREATE UNIQUE INDEX uk_oauth_accounts_user_provider
    ON auth.oauth_accounts (user_id, provider);

CREATE INDEX idx_oauth_accounts_user_id
    ON auth.oauth_accounts (user_id);

CREATE INDEX idx_oauth_accounts_provider_email
    ON auth.oauth_accounts (provider_email);

ALTER TABLE auth.oauth_accounts
    ADD CONSTRAINT fk_oauth_accounts_user
        FOREIGN KEY (user_id)
        REFERENCES auth.users(id)
        ON DELETE CASCADE;

ALTER TABLE auth.oauth_accounts
    ADD CONSTRAINT chk_oauth_accounts_provider
        CHECK (provider IN ('google', 'github'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS auth.oauth_accounts
    DROP CONSTRAINT IF EXISTS chk_oauth_accounts_provider;

ALTER TABLE IF EXISTS auth.oauth_accounts
    DROP CONSTRAINT IF EXISTS fk_oauth_accounts_user;

DROP INDEX IF EXISTS idx_oauth_accounts_provider_email;
DROP INDEX IF EXISTS idx_oauth_accounts_user_id;
DROP INDEX IF EXISTS uk_oauth_accounts_user_provider;
DROP INDEX IF EXISTS uk_oauth_accounts_provider_user;

DROP TABLE IF EXISTS auth.oauth_accounts;
-- +goose StatementEnd
