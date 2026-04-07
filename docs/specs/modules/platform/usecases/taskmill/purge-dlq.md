# Purge DLQ

Remove all dead letter queue tasks from a specified queue.

> **type**: user_action

> **operation-id**: `purge-dlq`

> **access**: POST /api/v1/platform/purge-dlq

> **actor**: user

> **permissions**: `taskmill:manage`

> **implementation**: [usecase.go](../../../../../../internal/modules/platform/usecase/taskmill/purgedlq/usecase.go)

## Input

```json
{
  "queue_name": "string"
}
```

- `queue_name` (string, required) - Name of the queue whose DLQ to purge

## Output

Empty response.

## Execute

- Remove all DLQ tasks from the specified queue

- Record audit log
