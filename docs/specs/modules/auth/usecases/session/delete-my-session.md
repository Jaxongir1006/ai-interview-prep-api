# Delete My Session

Allows an authenticated user to revoke one of their own sessions by ID.

> **type**: user_action

> **operation-id**: `delete-my-session`

> **access**: POST /api/v1/auth/delete-my-session

> **actor**: user (authenticated)

> **permissions**: -

> **implementation**: [usecase.go](../../../../../../internal/modules/auth/usecase/session/deletemysession/usecase.go)

## Input

```json
{
  "session_id": 1 // required, integer
}
```

## Output

Empty response.

## Execute

- Get user ID from authenticated user context

- Find target session by ID

- Verify target session belongs to the current user

- Start UOW

- Delete target session

- Record audit log

- Apply UOW
