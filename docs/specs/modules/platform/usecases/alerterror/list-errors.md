# List Errors

Retrieve application errors with filtering, search, sorting, and pagination.

> **type**: user_action
> **operation-id**: `list-errors`
> **access**: GET /api/v1/platform/list-errors
> **actor**: user
> **permissions**: `alert:view`
> **implementation**: [usecase.go](../../../../../../internal/modules/platform/usecase/alerterror/listerrors/usecase.go)

## Input

- `code`: string, optional — filter by error code
- `service`: string, optional — filter by service name
- `operation`: string, optional — filter by operation
- `alerted`: bool, optional — filter by alerted status
- `created_from`: string, optional — filter errors created after this timestamp (RFC3339)
- `created_to`: string, optional — filter errors created before this timestamp (RFC3339)
- `search`: string, optional — search across code, message, operation
- `sort`: string, optional — sort fields (allowed: created_at, code, service, operation). Default: created_at:desc
- `page_number`: int, optional — page number (default: 1)
- `page_size`: int, optional — page size (default: 20, max: 100)

## Output

```json
{
  "page_number": 1,
  "page_size": 20,
  "count": 150,
  "content": [
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
  ]
}
```

## Execute

- Normalize pagination params
- Parse sort options
- Build filter
- List errors with filters, search, sorting, and pagination
- Return paginated errors
