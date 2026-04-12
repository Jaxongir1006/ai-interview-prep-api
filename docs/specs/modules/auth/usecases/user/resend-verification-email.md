# Resend Verification Email

Sends a fresh email-verification link to a password-registered public user whose email is not verified yet.

> **type**: user_action

> **operation-id**: `resend-verification-email`

> **access**: POST /api/v1/auth/resend-verification-email

> **actor**: user (unauthenticated)

> **permissions**: -

> **implementation**: [usecase.go](../../../../../../internal/modules/auth/usecase/user/resendverificationemail/usecase.go)

## Input

```json
{
  "email": "user@example.com" // required, email format, max=255
}
```

## Output

Empty response.

## Execute

- Validate input

- Normalize email

- Find user by email

- Return empty response when no user exists for the email

- Check that user has password-based credentials

- Return empty response when user does not have password-based credentials

- Check that user is active

- Return empty response when user is not active

- Check that user is not already verified

- Return empty response when user is already verified

- Start UOW

- Expire previous unused email verification tokens for this user and email

- Create fresh one-time email verification token

- Apply UOW

- Send verification email with frontend verification URL and raw token

- Return empty response

## Security

- Return empty success response when no email is sent to avoid account enumeration
- Apply rate limiting by email and client IP before sending the message
- The raw token is never stored directly; only its hash is stored
