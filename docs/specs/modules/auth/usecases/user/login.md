# Login

Authenticates a public platform user with email and password, then creates a session with access and refresh tokens.

> **type**: user_action

> **operation-id**: `login`

> **access**: POST /api/v1/auth/login

> **actor**: user (unauthenticated)

> **permissions**: -

> **implementation**: [usecase.go](../../../../../../internal/modules/auth/usecase/user/login/usecase.go)

## Input

```json
{
  "email": "user@example.com", // required, email format, max=255
  "password": "string" // required, min=8, max=72
}
```

## Output

```json
{
  "access_token": "string",
  "access_token_expires_at": "2024-01-01T01:00:00Z",
  "refresh_token": "string",
  "refresh_token_expires_at": "2024-01-08T00:00:00Z"
}
```

## Execute

- Find user by email

- Check if user is active

- Verify password hash

- Start UOW

- Enforce max active sessions limit (delete least recently used sessions if exceeded)

- Create session record with tokens and meta info (IP, user_agent)

- Update user's last_login_at and last_active_at timestamps

- Record audit log

- Apply UOW

- Return session tokens
