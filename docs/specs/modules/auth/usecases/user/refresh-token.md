# Refresh Token

Generates new access and refresh tokens using a valid refresh token, extending the user's session.

> **type**: user_action

> **operation-id**: `refresh-token`

> **access**: POST /api/v1/auth/refresh-token

> **actor**: user (unauthenticated)

> **permissions**: -

> **implementation**: [usecase.go](../../../../../../internal/modules/auth/usecase/user/refreshtoken/usecase.go)

## Input

```json
{
  "refresh_token": "string" // required
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

- Find session by refresh token

- Check if refresh token is not expired

- Start UOW

- Generate new access and refresh tokens with updated expiry

- Update session record with new tokens and meta info (IP, user_agent)

- Update user's last_active_at timestamp

- Apply UOW

- Return new session tokens
