# Get Roles

Returns a filtered list of roles.

> **type**: user_action

> **operation-id**: `get-roles`

> **access**: GET /api/v1/auth/get-roles

> **actor**: user

> **permissions**: `auth:role:read`

> **implementation**: [usecase.go](../../../../../../internal/modules/auth/usecase/rbac/getroles/usecase.go)

## Input

- `page_number`: integer, optional, default 1
- `page_size`: integer, optional, default 20, max 100
- `id`: integer, optional
- `name`: string, optional
- `sort`: string, optional — sortable fields: name, created_at, updated_at. Default: created_at:desc

## Output

```json
{
  "page_number": 1,
  "page_size": 20,
  "count": 150,
  "content": [
    {
      "id": 1,
      "name": "string",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

## Execute

- Normalize pagination params

- List roles matching filter criteria

- Return paginated roles
