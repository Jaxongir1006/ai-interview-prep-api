# Update User

Updates an existing user's username and/or password.

> **type**: user_action

> **operation-id**: `update-user`

> **access**: POST /api/v1/auth/update-user

> **actor**: user

> **permissions**: `auth:user:manage`

> **implementation**: [usecase.go](../../../../../../internal/modules/auth/usecase/user/updateuser/usecase.go)

## Input

```json
{
  "id": "string", // required, UUID format
  "username": "string", // optional, min=3, max=50
  "password": "string" // optional, min=8
}
```

## Output

```json
{
  "id": "string",
  "username": "string",
  "is_active": true,
  "last_active_at": "2024-01-01T00:00:00Z", // nullable
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

## Execute

- Find user by ID

- If password provided, hash the new password

- Start UOW

- Update user fields

- Record audit log

- Apply UOW

- Return updated user
