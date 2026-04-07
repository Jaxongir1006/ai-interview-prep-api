# Get Action Logs

Returns a filtered list of action logs within a specified time range.

> **type**: user_action

> **operation-id**: `get-action-logs`

> **access**: GET /api/v1/audit/get-action-logs

> **actor**: user

> **permissions**: `audit:action-log:read`

> **implementation**: [usecase.go](../../../../../../internal/modules/audit/usecase/actionlog/getactionlogs/usecase.go)

## Input

- `from`: string, required, RFC3339 timestamp
- `to`: string, required, RFC3339 timestamp
- `module`: string, optional
- `operation_id`: string, optional
- `user_id`: string, optional
- `tags`: string[], optional, filter logs containing all specified tags
- `group_key`: string, optional, filter by business group identifier
- `cursor`: int64, optional
- `limit`: int, optional

## Output

```json
[
  {
    "id": 1,
    "user_id": "string", // nullable
    "module": "string",
    "operation_id": "string",
    "request_payload": {}, // nullable
    "tags": ["string"],
    "group_key": "string", // nullable
    "ip_address": "string",
    "user_agent": "string",
    "trace_id": "string",
    "created_at": "2026-01-01T00:00:00Z"
  }
]
```

## Execute

- Parse and validate time range (from < to)

- List action logs matching filter criteria

- Return action logs
