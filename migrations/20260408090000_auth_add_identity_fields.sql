-- +goose Up
-- +goose StatementBegin
ALTER TABLE auth.users
    ADD COLUMN IF NOT EXISTS email VARCHAR,
    ADD COLUMN IF NOT EXISTS phone_number VARCHAR,
    ADD COLUMN IF NOT EXISTS is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS last_login_at TIMESTAMPTZ;

CREATE UNIQUE INDEX IF NOT EXISTS uk_users_email
    ON auth.users (email)
    WHERE email IS NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS uk_users_phone_number
    ON auth.users (phone_number)
    WHERE phone_number IS NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS uk_users_phone_number;
DROP INDEX IF EXISTS uk_users_email;

ALTER TABLE auth.users
    DROP COLUMN IF EXISTS last_login_at,
    DROP COLUMN IF EXISTS is_verified,
    DROP COLUMN IF EXISTS phone_number,
    DROP COLUMN IF EXISTS email;
-- +goose StatementEnd
