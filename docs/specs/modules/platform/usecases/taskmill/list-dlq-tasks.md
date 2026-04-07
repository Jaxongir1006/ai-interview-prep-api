# List DLQ Tasks

Retrieve tasks from the dead letter queue with optional filtering.

> **type**: user_action

> **operation-id**: `list-dlq-tasks`

> **access**: GET /api/v1/platform/list-dlq-tasks

> **actor**: user

> **permissions**: `taskmill:view`

> **implementation**: [usecase.go](../../../../../../internal/modules/platform/usecase/taskmill/listdlqtasks/usecase.go)

## Input

- `queue_name` (string, optional) - Filter by queue name
- `operation_id` (string, optional) - Filter by operation ID
- `dlq_after` (time, optional) - Filter by DLQ timestamp after
- `dlq_before` (time, optional) - Filter by DLQ timestamp before
- `limit` (integer, optional) - Maximum number of results
- `offset` (integer, optional) - Pagination offset

## Output

```json
{
  "content": [
    {
      "id": 1,
      "queue_name": "string",
      "task_group_id": "string",
      "operation_id": "string",
      "payload": {},
      "priority": 0,
      "attempts": 3,
      "max_attempts": 3,
      "idempotency_key": "string",
      "created_at": "2024-01-01T00:00:00Z",
      "dlq_at": "2024-01-01T00:00:00Z",
      "dlq_reason": {}
    }
  ]
}
```

Note: `task_group_id` is nullable.

## Execute

- List DLQ tasks with filters
- Return DLQ tasks
