# Purge Queue

Remove all non-DLQ tasks from a specified queue.

> **type**: user_action

> **operation-id**: `purge-queue`

> **access**: POST /api/v1/platform/purge-queue

> **actor**: user

> **permissions**: `taskmill:manage`

> **implementation**: [usecase.go](../../../../../../internal/modules/platform/usecase/taskmill/purgequeue/usecase.go)

## Input

```json
{
  "queue_name": "string"
}
```

- `queue_name` (string, required) - Name of the queue to purge

## Output

Empty response.

## Execute

- Remove all non-DLQ tasks from the specified queue

- Record audit log
