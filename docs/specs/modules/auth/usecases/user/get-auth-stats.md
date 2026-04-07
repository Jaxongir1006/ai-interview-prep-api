# Get Auth Stats

Retrieve aggregated statistics for the auth module: total users, total roles, and active sessions.

> **type**: user_action

> **operation-id**: `get-auth-stats`

> **access**: GET /api/v1/auth/get-auth-stats

> **actor**: user

> **permissions**: `auth:user:read`

> **implementation**: [usecase.go](../../../../../../internal/modules/auth/usecase/user/getauthstats/usecase.go)

## Input

None.

## Output

```json
{
  "total_users": 0,
  "total_roles": 0,
  "active_sessions": 0
}
```

## Execute

- Count total users
- Count total roles
- Count active sessions (where refresh token has not expired)
- Return aggregated stats
