# Register

Registers a public platform user with email and password, initializes a minimal candidate profile, and sends an email-verification link.

> **type**: user_action

> **operation-id**: `register`

> **access**: POST /api/v1/auth/register

> **actor**: user (unauthenticated)

> **permissions**: -

> **implementation**: [usecase.go](../../../../../../internal/modules/auth/usecase/user/register/usecase.go)

## Input

```json
{
  "email": "user@example.com", // required, email format, max=255
  "password": "string", // required, min=8, max=72
  "full_name": "John Doe" // required, min=1, max=255
}
```

## Output

```json
{
  "email": "user@example.com",
  "verification_required": true
}
```

## Execute

- Validate input

- Check whether a user already exists with the same email
  - If the existing user is active, password-based, and not verified, create a fresh email verification token, send a new verification email, and return the normal success response
  - If the existing user is already verified, inactive, or OAuth-only, return `EMAIL_CONFLICT`

- Start UOW

- Create auth user with email-based credentials

- Set user `is_verified` to `false`

- Create minimal candidate profile for the new user using the provided full name

- Create one-time email verification token for the user's email

- Record audit log

- Apply UOW

- Send verification email with frontend verification URL and raw token

- Return normalized email with `verification_required = true`

## Email Verification

- Email/password registration requires the user to verify email ownership before password login succeeds
- The verification email links to the frontend verification page with a raw one-time token
- The raw token is never stored directly; only its hash is stored
- Verification is completed by `verify-email`
- OAuth login use cases do not send this email
- If email delivery fails after registration is committed, the user can request a fresh link through `resend-verification-email`
- Re-registering with the same unverified password-based email sends a fresh verification link without changing the stored password or profile data

## Errors

- Return `EMAIL_CONFLICT` when a verified, inactive, or OAuth-only user already exists with the same email
