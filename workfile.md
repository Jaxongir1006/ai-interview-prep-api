# Password Reset Backend Integration Plan

## Goal

Connect the existing `/password-reset` frontend route to real backend APIs and make the flow production-ready for future deployments.

The frontend currently has only a mock request flow. The backend needs to provide secure request and confirmation endpoints, email delivery, token lifecycle management, and stable error contracts that match the existing frontend API client conventions.

## Required User Flow

1. User opens `/password-reset`.
2. User enters their email address.
3. Frontend calls the password reset request endpoint.
4. Backend always returns a generic success response so account existence is not leaked.
5. If the email belongs to a local email/password account, backend sends a reset link.
6. User opens the reset link from email.
7. Frontend shows a new-password form using the token from the URL.
8. Frontend calls the password reset confirm endpoint.
9. Backend validates the token, updates the password, invalidates existing sessions/refresh tokens, and marks the token as used.
10. Frontend redirects the user to `/login` with a success message.

## Endpoints

### Request Password Reset

```http
POST /api/v1/auth/password-reset/request
Content-Type: application/json
```

Request:

```json
{
  "email": "john@example.com"
}
```

Success response:

```json
{
  "message": "If that email can receive password reset instructions, a reset link is on the way."
}
```

Important behavior:

- Always return `200 OK` for syntactically valid email requests, even if no account exists.
- Normalize email with trim + lowercase before lookup.
- Do not reveal whether the account exists.
- Do not send reset emails for OAuth-only accounts unless they also have a local password credential.
- Rate limit by email and source IP.
- If a valid unexpired reset token already exists, either reuse it or invalidate and create a new one. Prefer invalidating and creating a new token for simpler audit behavior.

Validation failure response:

```json
{
  "code": "VALIDATION_FAILED",
  "message": "Check the highlighted fields, then try again.",
  "fields": {
    "email": "Enter a valid email address."
  },
  "trace_id": "req_..."
}
```

### Confirm Password Reset

```http
POST /api/v1/auth/password-reset/confirm
Content-Type: application/json
```

Request:

```json
{
  "token": "raw-reset-token-from-email-link",
  "password": "newStrongPassword123"
}
```

Success response:

```json
{
  "message": "Your password has been reset. Sign in with your new password."
}
```

Failure responses:

Expired token:

```json
{
  "code": "PASSWORD_RESET_TOKEN_EXPIRED",
  "message": "This password reset link has expired. Request a new one.",
  "trace_id": "req_..."
}
```

Invalid or used token:

```json
{
  "code": "PASSWORD_RESET_TOKEN_INVALID",
  "message": "This password reset link is invalid or has already been used.",
  "trace_id": "req_..."
}
```

Weak password:

```json
{
  "code": "VALIDATION_FAILED",
  "message": "Check the highlighted fields, then try again.",
  "fields": {
    "password": "Password must be at least 8 characters."
  },
  "trace_id": "req_..."
}
```

## Token Requirements

- Generate a cryptographically secure random token with at least 32 bytes of entropy.
- Send only the raw token in the email link.
- Store only a hash of the token in the database.
- Use a short expiration window, recommended: 15 to 30 minutes.
- Token must be single-use.
- Store `created_at`, `expires_at`, `used_at`, and `requested_ip` where available.
- Invalidate older unused reset tokens for the same user when issuing a new one.
- On successful reset, invalidate all existing refresh tokens/sessions for that user.

Suggested table:

```sql
password_reset_tokens
- id
- user_id
- token_hash
- expires_at
- used_at
- requested_ip
- created_at
```

## Email Link Contract

The reset email should link back to the frontend app, not the backend API.

Recommended URL:

```text
{FRONTEND_URL}/password-reset?token={raw_token}
```

Backend deployment config should include:

- `FRONTEND_URL`
- `SMTP_HOST`
- `SMTP_PORT`
- `SMTP_USERNAME`
- `SMTP_PASSWORD`
- `EMAIL_FROM`
- Optional `PASSWORD_RESET_TOKEN_TTL_MINUTES`

For local development, `FRONTEND_URL` can point to the Vite dev server.

## Frontend Integration Expectations

The frontend API client already supports this error shape:

- `code`
- `message`
- `fields`
- `trace_id`

Add these API functions in `src/app/lib/api.ts` when backend is ready:

```ts
export type RequestPasswordResetRequest = {
  email: string;
};

export type ConfirmPasswordResetRequest = {
  token: string;
  password: string;
};

export async function requestPasswordReset(payload: RequestPasswordResetRequest) {
  return request<{ message: string }>(
    "/api/v1/auth/password-reset/request",
    {
      method: "POST",
      body: JSON.stringify(payload),
    },
    { suppressGlobalError: true },
  );
}

export async function confirmPasswordReset(payload: ConfirmPasswordResetRequest) {
  return request<{ message: string }>(
    "/api/v1/auth/password-reset/confirm",
    {
      method: "POST",
      body: JSON.stringify(payload),
    },
    { suppressGlobalError: true },
  );
}
```

Frontend page changes needed:

- If `/password-reset` has no `token`, show the email request form.
- If `/password-reset?token=...` exists, show a new password + confirm password form.
- Show field-level backend validation errors.
- Show generic success after request, regardless of account existence.
- Redirect to `/login` after successful confirm.
- Do not store the reset token in localStorage or sessionStorage.

## Security And Abuse Controls

- Rate limit reset requests by IP and normalized email.
- Add audit logs for reset requested, token used, token expired, and reset failed.
- Do not include raw tokens in logs.
- Do not reveal user existence in request endpoint responses.
- Use the same password hashing policy as registration.
- Reject passwords that exceed the backend password hashing limit if applicable.
- Consider rejecting password reuse if the backend stores password history.
- If the account is disabled, do not send a reset email and still return the generic success message.

## Deployment Readiness Checklist

- Backend has `FRONTEND_URL` configured per environment.
- Email provider credentials are configured in staging and production.
- Reset links use HTTPS outside local development.
- Token hash storage is migrated.
- Rate limiting is enabled.
- Existing refresh tokens are invalidated after reset.
- Logs include `trace_id` but never raw reset tokens.
- Frontend `VITE_API_URL` points to the backend environment.
- End-to-end smoke test covers request email and confirm password reset.

## Test Cases

- Valid email request returns generic success and sends email.
- Unknown email request returns the same generic success and sends no email.
- Invalid email returns `VALIDATION_FAILED` with `fields.email`.
- Valid token + valid password resets password.
- Expired token returns `PASSWORD_RESET_TOKEN_EXPIRED`.
- Used token returns `PASSWORD_RESET_TOKEN_INVALID`.
- Random token returns `PASSWORD_RESET_TOKEN_INVALID`.
- Weak password returns `VALIDATION_FAILED` with `fields.password`.
- Successful reset invalidates old refresh tokens.
- OAuth-only account request does not leak account state.
