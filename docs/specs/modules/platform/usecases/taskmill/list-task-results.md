# List Task Results

Retrieve completed task results with optional filtering.

> **type**: user_action

> **operation-id**: `list-task-results`

> **access**: GET /api/v1/platform/list-task-results

> **actor**: user

> **permissions**: `taskmill:view`

> **implementation**: [usecase.go](../../../../../../internal/modules/platform/usecase/taskmill/listtaskresults/usecase.go)

## Input

- `queue_name` (string, optional) - Filter by queue name
- `task_group_id` (string, optional) - Filter by task group ID
- `completed_after` (time, optional) - Filter by completion timestamp after
- `completed_before` (time, optional) - Filter by completion timestamp before
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
      "attempts": 1,
      "max_attempts": 3,
      "idempotency_key": "string",
      "scheduled_at": "2024-01-01T00:00:00Z",
      "created_at": "2024-01-01T00:00:00Z",
      "completed_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

Note: `task_group_id` is nullable.

## Execute

- List completed task results with filters
- Return task results
