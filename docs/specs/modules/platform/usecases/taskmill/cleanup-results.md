# Cleanup Results

Delete old task results matching the specified criteria.

> **type**: user_action

> **operation-id**: `cleanup-results`

> **access**: POST /api/v1/platform/cleanup-results

> **actor**: user

> **permissions**: `taskmill:manage`

> **implementation**: [usecase.go](../../../../../../internal/modules/platform/usecase/taskmill/cleanupresults/usecase.go)

## Input

```json
{
  "completed_before": "2024-01-01T00:00:00Z",
  "queue_name": "string"
}
```

- `completed_before` (time, required) - Delete results completed before this timestamp
- `queue_name` (string, optional) - Filter by queue name

## Output

```json
{
  "deleted_count": 0
}
```

## Execute

- Delete old task results matching criteria
- Record audit log
- Return count of deleted records
