-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS auth;

CREATE TABLE auth.users (
    id VARCHAR PRIMARY KEY,
    username VARCHAR UNIQUE,
    password_hash VARCHAR,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    last_active_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE auth.roles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE auth.role_permissions (
    id BIGSERIAL PRIMARY KEY,
    role_id BIGINT NOT NULL,
    permission VARCHAR NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_role_permissions_role_id ON auth.role_permissions (role_id);

CREATE TABLE auth.user_roles (
    id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR NOT NULL,
    role_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_user_roles_user_id ON auth.user_roles (user_id);

CREATE INDEX idx_user_roles_role_id ON auth.user_roles (role_id);

CREATE TABLE auth.user_permissions (
    id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR NOT NULL,
    permission VARCHAR NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_user_permissions_user_id ON auth.user_permissions (user_id);

CREATE TABLE auth.sessions (
    id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR NOT NULL,
    access_token VARCHAR NOT NULL,
    access_token_expires_at TIMESTAMPTZ NOT NULL,
    refresh_token VARCHAR NOT NULL,
    refresh_token_expires_at TIMESTAMPTZ NOT NULL,
    ip_address VARCHAR NOT NULL,
    user_agent VARCHAR NOT NULL,
    last_used_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_sessions_user_id ON auth.sessions (user_id);

CREATE INDEX idx_sessions_access_token ON auth.sessions (access_token);

CREATE INDEX idx_sessions_refresh_token ON auth.sessions (refresh_token);

CREATE INDEX idx_sessions_refresh_token_expires_at ON auth.sessions (refresh_token_expires_at);

ALTER TABLE auth.role_permissions ADD CONSTRAINT fk_role_permissions_role FOREIGN KEY (role_id) REFERENCES auth.roles (id) ON DELETE CASCADE;

ALTER TABLE auth.user_roles ADD CONSTRAINT fk_user_roles_user FOREIGN KEY (user_id) REFERENCES auth.users (id) ON DELETE CASCADE;

ALTER TABLE auth.user_roles ADD CONSTRAINT fk_user_roles_role FOREIGN KEY (role_id) REFERENCES auth.roles (id) ON DELETE CASCADE;

ALTER TABLE auth.user_permissions ADD CONSTRAINT fk_user_permissions_user FOREIGN KEY (user_id) REFERENCES auth.users (id) ON DELETE CASCADE;

ALTER TABLE auth.sessions ADD CONSTRAINT fk_sessions_user FOREIGN KEY (user_id) REFERENCES auth.users (id) ON DELETE CASCADE;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
-- Drop foreign keys first
ALTER TABLE IF EXISTS auth.sessions
DROP CONSTRAINT IF EXISTS fk_sessions_user;

ALTER TABLE IF EXISTS auth.user_permissions
DROP CONSTRAINT IF EXISTS fk_user_permissions_user;

ALTER TABLE IF EXISTS auth.user_roles
DROP CONSTRAINT IF EXISTS fk_user_roles_role;

ALTER TABLE IF EXISTS auth.user_roles
DROP CONSTRAINT IF EXISTS fk_user_roles_user;

ALTER TABLE IF EXISTS auth.role_permissions
DROP CONSTRAINT IF EXISTS fk_role_permissions_role;

-- Drop tables in reverse order
DROP TABLE IF EXISTS auth.sessions;

DROP TABLE IF EXISTS auth.user_permissions;

DROP TABLE IF EXISTS auth.user_roles;

DROP TABLE IF EXISTS auth.role_permissions;

DROP TABLE IF EXISTS auth.roles;

DROP TABLE IF EXISTS auth.users;

-- Drop schema
DROP SCHEMA IF EXISTS auth;

-- +goose StatementEnd
