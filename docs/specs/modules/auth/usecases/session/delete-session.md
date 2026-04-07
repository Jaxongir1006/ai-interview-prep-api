# Delete Session

Revokes a specific session by its ID. Used by administrators to force logout a user from a specific device.

> **type**: user_action

> **operation-id**: `delete-session`

> **access**: POST /api/v1/auth/delete-session

> **actor**: user

> **permissions**: `auth:session:manage`

> **implementation**: [usecase.go](../../../../../../internal/modules/auth/usecase/session/deletesession/usecase.go)

## Input

```json
{
  "session_id": 1 // required, integer
}
```

## Output

Empty response.

## Execute

- Find session by ID

- Start UOW

- Delete session

- Record audit log

- Apply UOW
