# Request Password Reset

Starts a password reset flow for a public email/password user by sending a one-time reset link when the email can receive one.

> **type**: user_action

> **operation-id**: `request-password-reset`

> **access**: POST /api/v1/auth/request-password-reset

> **actor**: user (unauthenticated)

> **permissions**: -

> **implementation**: [usecase.go](../../../../../../internal/modules/auth/usecase/user/requestpasswordreset/usecase.go)

## Input

```json
{
  "email": "user@example.com" // required, email format, max=255
}
```

## Output

```json
{
  "message": "If that email can receive password reset instructions, a reset link is on the way."
}
```

## Execute

- Validate input

- Normalize email

- Apply password reset request rate limit by normalized email and client IP

- Find user by email

- Return generic success response when no user exists for the email

- Check that user has password-based credentials

- Return generic success response when user does not have password-based credentials

- Check that user is active

- Return generic success response when user is not active

- Invalidate previous Redis-backed password reset token for this user and email

- Create fresh Redis-backed one-time password reset token for the user's current email with `auth.password_reset_token_ttl`

- Send password reset email with `auth.frontend_password_reset_url` and raw token

- Return generic success response

## Security

- Return the same success response for unknown, OAuth-only, inactive, and eligible accounts to avoid account enumeration
- The raw token is never stored directly; only its hash and reset metadata are stored in Redis
- The reset link points to the frontend password reset page and includes the raw token as a query parameter
- Do not include raw password reset tokens in logs, audit details, errors, or alerts
- Password reset token lifecycle does not produce audit records

## Errors

- Return `VALIDATION_FAILED` when email is missing or invalid
