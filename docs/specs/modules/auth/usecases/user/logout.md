# Logout

Terminates the current user session by deleting the session record.

> **type**: user_action

> **operation-id**: `logout`

> **access**: POST /api/v1/auth/logout

> **actor**: user (authenticated)

> **permissions**: -

> **implementation**: [usecase.go](../../../../../../internal/modules/auth/usecase/user/logout/usecase.go)

## Input

No input required. User identity is provided by auth middleware via the `Authorization` header.

## Output

Empty response.

## Execute

- Get session ID from authenticated user context

- Find session by ID

- Start UOW

- Delete session

- Record audit log

- Apply UOW
