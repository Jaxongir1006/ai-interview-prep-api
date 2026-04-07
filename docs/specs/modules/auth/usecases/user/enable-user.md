# Enable User

Reactivates a previously disabled user account, allowing them to log in again.

> **type**: user_action

> **operation-id**: `enable-user`

> **access**: POST /api/v1/auth/enable-user

> **actor**: user

> **permissions**: `auth:user:manage`

> **implementation**: [usecase.go](../../../../../../internal/modules/auth/usecase/user/enableuser/usecase.go)

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

- Check if user is not already active

- Start UOW

- Set user is_active to true

- Record audit log

- Apply UOW
