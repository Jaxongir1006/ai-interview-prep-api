# Get Error

Retrieve a single application error by ID with full details.

> **type**: user_action
> **operation-id**: `get-error`
> **access**: GET /api/v1/platform/get-error
> **actor**: user
> **permissions**: `alert:view`
> **implementation**: [usecase.go](../../../../../../internal/modules/platform/usecase/alerterror/geterror/usecase.go)

## Input

- `id`: string, required, uuid format — error record ID

## Output

{
  "id": "uuid",
  "code": "string",
  "message": "string",
  "details": { "key": "value" },
  "service": "string",
  "operation": "string",
  "created_at": "2024-01-01T00:00:00Z",
  "alerted": true
}

## Errors

- ERROR_NOT_FOUND (404) — error with given ID does not exist

## Execute

- Find error by ID
- Return error details
