# Admin Login

Authenticates a user with admin permissions and creates a session with access and refresh tokens.

> **type**: user_action

> **operation-id**: `admin-login`

> **access**: POST /api/v1/auth/admin-login

> **actor**: user (unauthenticated)

> **permissions**: -

> **implementation**: [usecase.go](../../../../../../internal/modules/auth/usecase/user/adminlogin/usecase.go)

## Input

```json
{
  "username": "string", // required
  "password": "string" // required
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

- Find user by username

- Check if user is active

- Verify password hash

- Start UOW

- Enforce max active sessions limit (delete least recently used sessions if exceeded)

- Create session record with tokens and meta info (IP, user_agent)

- Update user's last_active_at timestamp

- Record audit log

- Apply UOW

- Return session tokens
