# Confirm Password Reset

Completes a password reset using a one-time token from email, updates the user's password, and revokes existing sessions.

> **type**: user_action

> **operation-id**: `confirm-password-reset`

> **access**: POST /api/v1/auth/confirm-password-reset

> **actor**: user (unauthenticated)

> **permissions**: -

> **implementation**: [usecase.go](../../../../../../internal/modules/auth/usecase/user/confirmpasswordreset/usecase.go)

## Input

```json
{
  "token": "string", // required, raw token from password reset link
  "password": "string" // required, min=8, max=72
}
```

## Output

```json
{
  "message": "Your password has been reset. Sign in with your new password."
}
```

## Execute

- Validate input

- Hash the raw token

- Consume Redis-backed password reset token by token hash

- Find user by token user ID

- Check that user is active

- Check that user has password-based credentials

- Check that token email matches the user's current email

- Hash the new password

- Start UOW

- Update user's password_hash

- Mark user as verified

- Delete all sessions for the user

- Apply UOW

- Return success response

## Security

- The raw token is never stored directly; only its hash and reset metadata are stored in Redis
- A successful reset verifies the user's current email because the raw token was delivered to that email address
- Existing sessions are deleted so previous access and refresh tokens cannot continue to be used after the reset
- Do not include raw password reset tokens in logs, audit details, errors, or alerts
- Password reset token lifecycle does not produce audit records

## Errors

- Return `PASSWORD_RESET_TOKEN_INVALID` when token is not found, already used, expired from Redis, belongs to an inactive user, belongs to a user without password-based credentials, or no longer matches the user's current email
- Return `VALIDATION_FAILED` when token is missing or password does not satisfy the password policy
