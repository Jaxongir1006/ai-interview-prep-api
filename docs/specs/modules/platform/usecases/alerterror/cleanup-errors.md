# Cleanup Errors

Delete old error records created before a specified timestamp.

> **type**: user_action
> **operation-id**: `cleanup-errors`
> **access**: POST /api/v1/platform/cleanup-errors
> **actor**: user
> **permissions**: `alert:manage`
> **implementation**: [usecase.go](../../../../../../internal/modules/platform/usecase/alerterror/cleanuperrors/usecase.go)

## Input

{
  "created_before": "2024-01-01T00:00:00Z"  // required, RFC3339
}

## Output

{
  "deleted_count": 5230
}

## Execute

- Delete errors older than the specified timestamp
- Record audit log
- Return count of deleted records
