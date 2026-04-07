# Disable User

Deactivates a user account and terminates all their active sessions, preventing further access.

> **type**: user_action

> **operation-id**: `disable-user`

> **access**: POST /api/v1/auth/disable-user

> **actor**: user

> **permissions**: `auth:user:manage`

> **implementation**: [usecase.go](../../../../../../internal/modules/auth/usecase/user/disableuser/usecase.go)

## Input

```json
{
  "id": "string" // required, UUID format
}
```

## Output

Empty response.

## Execute

- Find user by ID

- Check if user is not already disabled

- Start UOW

- Set user is_active to false

- Delete all user sessions (force logout)

- Record audit log

- Apply UOW
