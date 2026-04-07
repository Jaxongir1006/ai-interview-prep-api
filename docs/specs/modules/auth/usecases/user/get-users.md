# Get Users

Returns a filtered list of user accounts for administrative purposes.

> **type**: user_action

> **operation-id**: `get-users`

> **access**: GET /api/v1/auth/get-users

> **actor**: user

> **permissions**: `auth:user:read`

> **implementation**: [usecase.go](../../../../../../internal/modules/auth/usecase/user/getusers/usecase.go)

## Input

- `page_number`: integer, optional, default 1
- `page_size`: integer, optional, default 20, max 100
- `id`: string, optional, UUID format
- `username`: string, optional
- `is_active`: bool, optional
- `sort`: string, optional — sortable fields: username, created_at, updated_at. Default: created_at:desc

## Output

```json
{
  "page_number": 1,
  "page_size": 20,
  "count": 150,
  "content": [
    {
      "id": "string",
      "username": "string",
      "is_active": true,
      "last_active_at": "2024-01-01T00:00:00Z", // nullable
      "roles": ["string"],
      "direct_permissions": ["string"],
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

## Execute

- Normalize pagination params

- List users matching filter criteria

- Batch-fetch role assignments and role names for all users

- Batch-fetch direct permissions for all users

- Return users with roles and permissions
