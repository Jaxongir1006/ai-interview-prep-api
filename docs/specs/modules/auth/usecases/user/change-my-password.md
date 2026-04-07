# Change My Password

Allows authenticated users to change their own password by providing the current password and a new one.

> **type**: user_action

> **operation-id**: `change-my-password`

> **access**: POST /api/v1/auth/change-my-password

> **actor**: user (authenticated)

> **permissions**: -

> **implementation**: [usecase.go](../../../../../../internal/modules/auth/usecase/user/changemypassword/usecase.go)

## Input

```json
{
  "current_password": "string", // required
  "new_password": "string" // required, min=8
}
```

## Output

Empty response.

## Execute

- Get user context from authenticated session

- Find user by ID

- Verify current password matches stored hash
  - If mismatch -> INCORRECT_CREDENTIALS

- Hash new password

- Start UOW

- Update user's password_hash

- Record audit log

- Apply UOW
