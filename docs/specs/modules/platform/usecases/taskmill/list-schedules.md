# List Schedules

Retrieve all cron schedules, optionally filtered by queue name.

> **type**: user_action

> **operation-id**: `list-schedules`

> **access**: GET /api/v1/platform/list-schedules

> **actor**: user

> **permissions**: `taskmill:view`

> **implementation**: [usecase.go](../../../../../../internal/modules/platform/usecase/taskmill/listschedules/usecase.go)

## Input

- `queue_name` (string, optional) - Filter by queue name

## Output

```json
{
  "content": [
    {
      "operation_id": "string",
      "queue_name": "string",
      "cron_pattern": "string",
      "next_run_at": "2024-01-01T00:00:00Z",
      "last_run_at": "2024-01-01T00:00:00Z",
      "last_run_status": "string",
      "last_error": "string",
      "run_count": 0
    }
  ]
}
```

Note: `last_run_at`, `last_run_status`, and `last_error` are nullable.

## Execute

- List cron schedules, optionally filtered by queue
- Return schedules
