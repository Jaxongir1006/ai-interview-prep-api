# Requeue From DLQ

Move a task from the dead letter queue back to the main queue for retry.

> **type**: user_action

> **operation-id**: `requeue-from-dlq`

> **access**: POST /api/v1/platform/requeue-from-dlq

> **actor**: user

> **permissions**: `taskmill:manage`

> **implementation**: [usecase.go](../../../../../../internal/modules/platform/usecase/taskmill/requeuefromdlq/usecase.go)

## Input

```json
{
  "task_id": 1
}
```

- `task_id` (integer, required) - ID of the DLQ task to requeue

## Output

Empty response.

## Execute

- Move DLQ task back to the main queue for retry

- Record audit log
