# API General

Standard response structures for all API endpoints.

## Response Headers

Every response includes an `X-Trace-ID` header with the request trace ID for debugging and log correlation.

## CORS

The HTTP API supports configurable CORS at the shared server layer. Browser clients must send requests from an origin listed in `cors.allow_origins`.

Local development allows these origins by default:

```json
[
  "http://localhost:*",
  "http://127.0.0.1:*",
  "http://localhost:3000",
  "http://localhost:5173",
  "http://127.0.0.1:3000",
  "http://127.0.0.1:5173"
]
```

Local-only wildcard entries allow browser dev servers on any port for `localhost` and `127.0.0.1`.

Preflight requests use `OPTIONS` and return `204 No Content` when the origin and requested method are allowed. The default local policy allows `GET`, `POST`, and `OPTIONS`, accepts `Origin`, `Content-Type`, `Accept`, `Authorization`, and `X-Requested-With`, exposes `X-Trace-ID`, and supports credentials for explicit allowed origins.

## List Response

All list endpoints wrap items in a `content` field:

```json
{
  "content": []
}
```

## Paginated Response

Paginated endpoints extend the list response with pagination metadata:

```json
{
  "page_number": 1,
  "page_size": 20,
  "count": 150,
  "content": []
}
```

| Field         | Type  | Description                     |
| ------------- | ----- | ------------------------------- |
| `page_number` | int   | Current page number (1-indexed) |
| `page_size`   | int   | Items per page (max 100)        |
| `count`       | int   | Total number of matching items  |
| `content`     | array | Array of items for current page |

## Error Response

```json
{
  "trace_id": "string",
  "error": {
    "code": "string",
    "message": "string",
    "cause": "string",
    "fields": {
      "field_name": "error description"
    },
    "trace": "string",
    "details": {
      "key": "value"
    }
  }
}
```

| Field           | Type              | Description                                            |
| --------------- | ----------------- | ------------------------------------------------------ |
| `trace_id`      | string            | Request trace ID for debugging, same with `X-Trace-ID` |
| `error.code`    | string            | Machine-readable error code (e.g., `ADMIN_NOT_FOUND`)  |
| `error.message` | string            | Translated user-facing message for the error code      |
| `error.cause`   | string            | Original technical error description                   |
| `error.trace`   | string            | Error stack trace (hidden in production)               |
| `error.fields`  | map[string]string | Field-level validation errors (optional)               |
| `error.details` | map[string]any    | Additional context (hidden in production)              |

### Validation Error Example

```json
{
  "trace_id": "abc-123",
  "error": {
    "code": "VALIDATION_FAILED",
    "message": "Введённые данные некорректны",
    "cause": "Validation failed. See fields for details.",
    "fields": {
      "username": "This field is required",
      "email": "Invalid email format",
      "password": "Must be at least 8 characters"
    }
  }
}
```

### Business Error Example

```json
{
  "trace_id": "abc-123",
  "error": {
    "code": "INCORRECT_CREDENTIALS",
    "message": "Неверное имя пользователя или пароль",
    "cause": "username or password is incorrect"
  }
}
```
