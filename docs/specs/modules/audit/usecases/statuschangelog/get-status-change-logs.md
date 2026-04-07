# Get Status Change Logs

Returns a filtered list of entity status change logs within a specified time range.

> **type**: user_action

> **operation-id**: `get-status-change-logs`

> **access**: GET /api/v1/audit/get-status-change-logs

> **actor**: user

> **permissions**: `audit:status-change-log:read`

> **implementation**: [usecase.go](../../../../../../internal/modules/audit/usecase/statuschangelog/getstatuschangelogs/usecase.go)

## Input

- `from`: string, required, RFC3339 timestamp format
- `to`: string, required, RFC3339 timestamp format
- `entity_type`: string, optional
- `entity_id`: string, optional
- `action_log_id`: int64, optional
- `cursor`: int64, optional
- `limit`: int, optional

## Output

```json
[
  {
    "id": 1,
    "action_log_id": 1,
    "entity_type": "string",
    "entity_id": "string",
    "status": "string",
    "trace_id": "string",
    "created_at": "2024-01-01T00:00:00Z"
  }
]
```

## Execute

- Parse and validate time range parameters

- List status change logs matching filter criteria

- Return status change logs
