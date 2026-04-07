# List Queues

Retrieve all distinct queue names from the task system.

> **type**: user_action

> **operation-id**: `list-queues`

> **access**: GET /api/v1/platform/list-queues

> **actor**: user

> **permissions**: `taskmill:view`

> **implementation**: [usecase.go](../../../../../../internal/modules/platform/usecase/taskmill/listqueues/usecase.go)

## Input

- None

## Output

```json
{
  "content": [
    "queue1",
    "queue2"
  ]
}
```

## Execute

- List all distinct queue names
- Return queue names
