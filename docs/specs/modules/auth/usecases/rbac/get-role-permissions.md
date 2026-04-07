# Get Role Permissions

Returns all permissions assigned to a specific role.

> **type**: user_action

> **operation-id**: `get-role-permissions`

> **access**: GET /api/v1/auth/get-role-permissions

> **actor**: user

> **permissions**: `auth:role:read`

> **implementation**: [usecase.go](../../../../../../internal/modules/auth/usecase/rbac/getrolepermissions/usecase.go)

## Input

- `role_id`: integer, required

## Output

```json
{
  "content": [
    {
      "id": 1,
      "role_id": 1,
      "permission": "string",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

## Execute

- Find role by ID

- List permissions for the role

- Return permissions
