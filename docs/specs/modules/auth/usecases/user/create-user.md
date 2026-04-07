# Create User

Creates a new user account with the provided credentials.

> **type**: user_action

> **operation-id**: `create-user`

> **access**: POST /api/v1/auth/create-user

> **actor**: user

> **permissions**: `auth:user:manage`

> **implementation**: [usecase.go](../../../../../../internal/modules/auth/usecase/user/createuser/usecase.go)

## Input

```json
{
  "username": "string", // required, min=3, max=50
  "password": "string" // required, min=8
}
```

## Output

```json
{
  "id": "string",
  "username": "string",
  "is_active": true,
  "last_active_at": null,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

## Execute

- Hash the password

- Start UOW

- Create user

- Record audit log

- Apply UOW

- Return created user
