# Delete User Sessions

Revokes all active sessions for a specific user, effectively forcing a complete logout from all devices.

> **type**: user_action

> **operation-id**: `delete-user-sessions`

> **access**: POST /api/v1/auth/delete-user-sessions

> **actor**: user

> **permissions**: `auth:session:manage`

> **implementation**: [usecase.go](../../../../../../internal/modules/auth/usecase/session/deleteusersessions/usecase.go)

## Input

```json
{
  "user_id": "string" // required, UUID format
}
```

## Output

Empty response.

## Execute

- Find user by ID

- Start UOW

- Delete all sessions for the user

- Record audit log

- Apply UOW
