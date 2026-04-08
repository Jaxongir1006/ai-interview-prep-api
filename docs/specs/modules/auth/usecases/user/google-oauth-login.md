# Google OAuth Login

Authenticates a public platform user with Google OAuth and creates or reuses the linked account, candidate profile, and session.

> **type**: user_action

> **operation-id**: `google-oauth-login`

> **access**: POST /api/v1/auth/google-oauth-login

> **actor**: user (unauthenticated)

> **permissions**: -

> **implementation**: [usecase.go](../../../../../../internal/modules/auth/usecase/user/googleoauthlogin/usecase.go)

## Input

```json
{
  "id_token": "string" // required, Google-issued ID token from the client OAuth flow
}
```

## Output

```json
{
  "access_token": "string",
  "access_token_expires_at": "2024-01-01T01:00:00Z",
  "refresh_token": "string",
  "refresh_token_expires_at": "2024-01-08T00:00:00Z",
  "is_new_user": true
}
```

## Execute

- Validate Google token with Google identity metadata

- Extract provider user ID, email, verification status, and full name from the provider payload

- Find linked OAuth account by provider and provider user ID

- Find user by email when no linked OAuth account exists

- Start UOW

- Create auth user when no matching user exists

- Create linked Google OAuth account when it does not exist

- Create minimal candidate profile when a new public user is created

- Enforce max active sessions limit (delete least recently used sessions if exceeded)

- Create session record with tokens and meta info (IP, user_agent)

- Update user's last_login_at and last_active_at timestamps

- Record audit log

- Apply UOW

- Return session tokens and whether the user was newly created
