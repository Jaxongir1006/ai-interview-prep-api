# Get Error Stats

Retrieve aggregated error statistics grouped by service, operation, and code.

> **type**: user_action
> **operation-id**: `get-error-stats`
> **access**: GET /api/v1/platform/get-error-stats
> **actor**: user
> **permissions**: `alert:view`
> **implementation**: [usecase.go](../../../../../../internal/modules/platform/usecase/alerterror/geterrorstats/usecase.go)

## Input

- `created_from`: time, optional — count errors created after this timestamp
- `created_to`: time, optional — count errors created before this timestamp

## Output

{
  "total_count": 1523,
  "by_service": [
    { "service": "my-app", "count": 1200 }
  ],
  "by_operation": [
    { "operation": "admin-login", "count": 450 }
  ],
  "by_code": [
    { "code": "INTERNAL_ERROR", "count": 800 }
  ]
}

## Execute

- Get error statistics for the given time range
- Return aggregated stats
