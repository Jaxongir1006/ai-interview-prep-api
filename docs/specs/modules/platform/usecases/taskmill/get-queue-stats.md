# Get Queue Stats

Retrieve statistics for a specified queue.

> **type**: user_action

> **operation-id**: `get-queue-stats`

> **access**: GET /api/v1/platform/get-queue-stats

> **actor**: user

> **permissions**: `taskmill:view`

> **implementation**: [usecase.go](../../../../../../internal/modules/platform/usecase/taskmill/getqueuestats/usecase.go)

## Input

- `queue_name` (string, required) - Name of the queue

## Output

```json
{
  "queue_name": "string",
  "total": 0,
  "available": 0,
  "in_flight": 0,
  "scheduled": 0,
  "in_dlq": 0,
  "oldest_task": "2024-01-01T00:00:00Z",
  "avg_attempts": 0.0,
  "p95_attempts": 0.0
}
```

Note: `oldest_task` is nullable.

## Execute

- Get stats for the specified queue
- Return queue statistics
