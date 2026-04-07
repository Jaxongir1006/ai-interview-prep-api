# Trigger Schedule

Manually trigger a scheduled task for immediate execution.

> **type**: user_action

> **operation-id**: `trigger-schedule`

> **access**: POST /api/v1/platform/trigger-schedule

> **actor**: user

> **permissions**: `taskmill:manage`

> **implementation**: [usecase.go](../../../../../../internal/modules/platform/usecase/taskmill/triggerschedule/usecase.go)

## Input

```json
{
  "operation_id": "string"
}
```

- `operation_id` (string, required) - Operation ID of the schedule to trigger

## Output

Empty response.

## Execute

- Verify schedule exists by operation ID
- Enqueue task to the schedule's queue for immediate execution
- Record audit log
