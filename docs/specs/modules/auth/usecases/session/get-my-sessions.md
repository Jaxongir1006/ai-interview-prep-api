# Get My Sessions

Returns the list of active sessions for the current authenticated user.

> **type**: user_action

> **operation-id**: `get-my-sessions`

> **access**: GET /api/v1/auth/get-my-sessions

> **actor**: user (authenticated)

> **permissions**: -

> **implementation**: [usecase.go](../../../../../../internal/modules/auth/usecase/session/getmysessions/usecase.go)

## Input

No input required. User identity is provided by auth middleware via the `Authorization` header.

## Output

```json
{
  "content": [
    {
      "id": 1,
      "user_id": "string",
      "access_token_expires_at": "2024-01-01T00:00:00Z",
      "refresh_token_expires_at": "2024-01-01T00:00:00Z",
      "ip_address": "string",
      "user_agent": "string",
      "last_used_at": "2024-01-01T00:00:00Z",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

## Execute

- Get user ID from authenticated user context

- List all sessions for the user

- Return sessions
