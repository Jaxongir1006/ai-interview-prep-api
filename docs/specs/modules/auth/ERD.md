# Auth Module ERD

```mermaid
erDiagram
    users {
        VARCHAR id PK "UUID-formatted string identifier"
        VARCHAR username UK "nullable, admin login identifier"
        VARCHAR email UK "nullable, public user login identity"
        VARCHAR phone_number UK "nullable, public user contact number"
        VARCHAR password_hash "nullable, supports future external auth providers"
        BOOLEAN is_verified "whether the user's primary public identity is verified"
        BOOLEAN is_active
        TIMESTAMPTZ last_login_at "updated on successful interactive login"
        TIMESTAMPTZ last_active_at "updated on login and token refresh"
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }

    oauth_accounts {
        BIGSERIAL id PK
        VARCHAR user_id FK
        VARCHAR provider "google or github"
        VARCHAR provider_user_id "provider-side stable user identifier"
        VARCHAR provider_email "nullable, email returned by provider"
        TIMESTAMPTZ last_login_at "updated on successful OAuth login"
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }

    roles {
        BIGSERIAL id PK
        VARCHAR name UK
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }

    role_permissions {
        BIGSERIAL id PK
        BIGINT role_id FK
        VARCHAR permission
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }

    user_roles {
        BIGSERIAL id PK
        VARCHAR user_id FK
        BIGINT role_id FK
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }

    user_permissions {
        BIGSERIAL id PK
        VARCHAR user_id FK
        VARCHAR permission
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }

    sessions {
        BIGSERIAL id PK
        VARCHAR user_id FK
        VARCHAR access_token
        TIMESTAMPTZ access_token_expires_at
        VARCHAR refresh_token
        TIMESTAMPTZ refresh_token_expires_at
        VARCHAR ip_address
        VARCHAR user_agent
        TIMESTAMPTZ last_used_at
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }

    users ||--o{ oauth_accounts : "linked to"
    users ||--o{ user_roles : "has"
    users ||--o{ user_permissions : "has"
    users ||--o{ sessions : "has"
    roles ||--o{ role_permissions : "has"
    roles ||--o{ user_roles : "assigned via"
```

## Notes

- `username`, `email`, and `phone_number` are nullable to support different actor types in the same `users` table
- Admin accounts authenticate with `username + password`; public users authenticate with `email + password`
- Public users may alternatively authenticate through linked `oauth_accounts` for Google and GitHub
- `oauth_accounts` should enforce uniqueness for `(provider, provider_user_id)`
- Email verification token state is stored in Redis, not PostgreSQL
- Password reset token state is stored in Redis, not PostgreSQL
- Successful password reset deletes all existing sessions for the user
- Business-profile and interview-preparation data for public users should live in a separate module such as `candidate`
- When implementing the migration, follow [Migration Guideline](../../../guidelines/13_db_migrations.md): create tables first, create indexes second, then add foreign keys and check constraints with `ALTER TABLE`

## Redis Token State

Email verification and password reset tokens are ephemeral auth state and are not part of the PostgreSQL ERD.

Redis stores only token hashes and metadata needed to complete the flow:

```json
{
  "user_id": "string",
  "email": "user@example.com",
  "expires_at": "2026-04-23T12:00:00Z"
}
```

Expected key ownership:

- `auth:email_verification:token:{token_hash}` stores email verification token metadata with `auth.email_verification_token_ttl`
- `auth:email_verification:user:{user_id}:{email}` points to the current email verification token hash so older tokens can be invalidated
- `auth:password_reset:token:{token_hash}` stores password reset token metadata with `auth.password_reset_token_ttl`
- `auth:password_reset:user:{user_id}:{email}` points to the current password reset token hash so older tokens can be invalidated

Token confirmation should consume token keys atomically so a token cannot be reused.
