# Get User Roles

Returns all roles assigned to a specific user.

> **type**: user_action

> **operation-id**: `get-user-roles`

> **access**: GET /api/v1/auth/get-user-roles

> **actor**: user

> **permissions**: `auth:access:read`

> **implementation**: [usecase.go](../../../../../../internal/modules/auth/usecase/rbac/getuserroles/usecase.go)

## Input

- `user_id`: string, required, UUID format

## Output

```json
{
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

- Find user by ID

- List role assignments for the user

- Fetch role details for each assignment

- Return roles
