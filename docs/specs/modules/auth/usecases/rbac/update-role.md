# Update Role

Updates an existing role's name.

> **type**: user_action

> **operation-id**: `update-role`

> **access**: POST /api/v1/auth/update-role

> **actor**: user

> **permissions**: `auth:role:manage`

> **implementation**: [usecase.go](../../../../../../internal/modules/auth/usecase/rbac/updaterole/usecase.go)

## Input

```json
{
  "id": 1, // required, integer
  "name": "string" // required, min=3, max=50
}
```

## Output

```json
{
  "id": 1,
  "name": "string",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

## Execute

- Find role by ID

- Start UOW

- Update role name

- Record audit log

- Apply UOW

- Return updated role
