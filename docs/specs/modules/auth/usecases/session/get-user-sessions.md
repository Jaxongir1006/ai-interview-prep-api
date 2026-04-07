# Get User Sessions

Returns a list of active sessions for a specific user. Used by administrators to monitor user sessions.

> **type**: user_action

> **operation-id**: `get-user-sessions`

> **access**: GET /api/v1/auth/get-user-sessions

> **actor**: user

> **permissions**: `auth:session:read`

> **implementation**: [usecase.go](../../../../../../internal/modules/auth/usecase/session/getusersessions/usecase.go)

## Input

- `page_number`: integer, optional, default 1
- `page_size`: integer, optional, default 20, max 100
- `user_id`: string, required, UUID format
- `sort`: string, optional — sortable fields: last_used_at, created_at. Default: created_at:desc

## Output

```json
{
  "page_number": 1,
  "page_size": 20,
  "count": 150,
  "content": [
    {
      "id": 1,
      "user_id": "string",
      "ip_address": "string",
      "user_agent": "string",
      "last_used_at": "2024-01-01T00:00:00Z",
      "created_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

## Execute

- Normalize pagination params

- Find user by ID

- List sessions for the user

- Return paginated sessions
