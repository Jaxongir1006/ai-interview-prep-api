# Create Role

Creates a new role that can be assigned permissions and associated with users.

> **type**: user_action

> **operation-id**: `create-role`

> **access**: POST /api/v1/auth/create-role

> **actor**: user

> **permissions**: `auth:role:manage`

> **implementation**: [usecase.go](../../../../../../internal/modules/auth/usecase/rbac/createrole/usecase.go)

## Input

```json
{
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

- Start UOW

- Create role

- Record audit log

- Apply UOW

- Return created role
