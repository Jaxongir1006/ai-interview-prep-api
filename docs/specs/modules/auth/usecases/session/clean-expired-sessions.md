# Clean Expired Sessions

Periodically removes sessions from the database where the refresh token has expired, preventing database bloat from accumulated expired session records.

> **type**: async_task

> **operation-id**: `clean-expired-sessions`

> **implementation**: [usecase.go](../../../../../../internal/modules/auth/usecase/session/cleanexpiredsessions/usecase.go)

## Task payload

```json
{}
```

## Execute

- Delete all sessions where refresh token has expired

## Idempotency

The task is inherently idempotent.

## Schedule

Runs hourly via cron pattern: `0 * * * *` (at minute 0 of every hour).
