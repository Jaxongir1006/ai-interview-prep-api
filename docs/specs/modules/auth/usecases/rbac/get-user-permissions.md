# Get User Permissions

Returns all direct permissions assigned to a specific user. Does not include permissions inherited from roles.

> **type**: user_action

> **operation-id**: `get-user-permissions`

> **access**: GET /api/v1/auth/get-user-permissions

> **actor**: user

> **permissions**: `auth:access:read`

> **implementation**: [usecase.go](../../../../../../internal/modules/auth/usecase/rbac/getuserpermissions/usecase.go)

## Input

- `user_id`: string, required, UUID format

## Output

```json
{
  "content": [
    {
      "id": 1,
      "user_id": "string",
      "permission": "string",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

## Execute

- Find user by ID

- List direct permissions for the user

- Return permissions
