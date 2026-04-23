-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS auth.email_verification_tokens
    DROP CONSTRAINT IF EXISTS fk_email_verification_tokens_user;

DROP INDEX IF EXISTS idx_email_verification_tokens_expires_at;
DROP INDEX IF EXISTS idx_email_verification_tokens_user_email;
DROP INDEX IF EXISTS uk_email_verification_tokens_token_hash;

DROP TABLE IF EXISTS auth.email_verification_tokens;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE TABLE auth.email_verification_tokens (
    id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR NOT NULL,
    email VARCHAR NOT NULL,
    token_hash VARCHAR NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX uk_email_verification_tokens_token_hash
    ON auth.email_verification_tokens (token_hash);

CREATE INDEX idx_email_verification_tokens_user_email
    ON auth.email_verification_tokens (user_id, email);

CREATE INDEX idx_email_verification_tokens_expires_at
    ON auth.email_verification_tokens (expires_at);

ALTER TABLE auth.email_verification_tokens
    ADD CONSTRAINT fk_email_verification_tokens_user
    FOREIGN KEY (user_id) REFERENCES auth.users(id)
    ON DELETE CASCADE;
-- +goose StatementEnd
