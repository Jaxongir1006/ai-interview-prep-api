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
  "user": {
    "id": "string",
    "email": "user@example.com",
    "phone_number": null,
    "is_verified": false,
    "is_active": true,
    "last_login_at": null,
    "last_active_at": null,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "profile": {
    "id": 1,
    "user_id": "string",
    "full_name": "John Doe",
    "bio": null,
    "location": null,
    "target_role": null,
    "experience_level": null,
    "interview_goal_per_week": 3,
    "preferred_topics": [],
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "verification_required": true
}
```

## Execute

- Validate input

- Check that no user already exists with the same email
  - If a user already exists with the same email, return `EMAIL_ALREADY_EXISTS`
  - If the existing user is not verified, the client should offer `resend-verification-email`

- Start UOW

- Create auth user with email-based credentials

- Set user `is_verified` to `false`

- Create minimal candidate profile for the new user using the provided full name

- Create one-time email verification token for the user's email

- Record audit log

- Apply UOW

- Send verification email with frontend verification URL and raw token

- Return created user and profile data with `verification_required = true`

## Email Verification

- Email/password registration requires the user to verify email ownership before password login succeeds
- The verification email links to the frontend verification page with a raw one-time token
- The raw token is never stored directly; only its hash is stored
- Verification is completed by `verify-email`
- OAuth login use cases do not send this email
- If email delivery fails after registration is committed, the user can request a fresh link through `resend-verification-email`

## Errors

- Return `EMAIL_ALREADY_EXISTS` when a user already exists with the same email
- When `EMAIL_ALREADY_EXISTS` is returned for an unverified password-registered user, the client should guide the user to `resend-verification-email`
