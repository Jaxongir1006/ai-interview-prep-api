# Verify Email

Verifies ownership of a password-registered public user's email address using a one-time token from the verification email.

> **type**: user_action

> **operation-id**: `verify-email`

> **access**: POST /api/v1/auth/verify-email

> **actor**: user (unauthenticated)

> **permissions**: -

> **implementation**: [usecase.go](../../../../../../internal/modules/auth/usecase/user/verifyemail/usecase.go)

## Input

```json
{
  "token": "string" // required, raw token from verification link
}
```

## Output

```json
{
  "user_id": "string",
  "email": "user@example.com",
  "is_verified": true,
  "access_token": "string",
  "access_token_expires_at": "2026-04-13T11:00:00Z",
  "refresh_token": "string",
  "refresh_token_expires_at": "2026-04-20T10:00:00Z",
  "onboarding_required": true
}
```

## Execute

- Validate input

- Hash the raw token

- Find unused email verification token by token hash

- Check that token is not expired

- Find user by token user ID

- Check that token email matches the user's current email

- Start UOW

- Mark email verification token as used

- Mark user as verified

- Create authenticated session

- Record audit log

- Apply UOW

- Resolve whether onboarding is still required

- Return verified user identity data, authenticated session tokens, and onboarding requirement

## Errors

- Return `EMAIL_VERIFICATION_TOKEN_INVALID` when token is not found or already used
- Return `EMAIL_VERIFICATION_TOKEN_EXPIRED` when token is expired
- Return `EMAIL_VERIFICATION_EMAIL_MISMATCH` when token email does not match the user's current email
