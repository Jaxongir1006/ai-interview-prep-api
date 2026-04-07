# {Use Case Name}

{One-two sentence summary of what this does and why.}

> **type**: user_action

> **operation-id**: `{operation-id}`

> **access**: {GET or POST} {path}

> **actor**: {one of user, admin, service_acc}

> **permissions**: `{permissions list}`

> **implementation**: [usecase.go](../../../../internal/modules/{module}/usecase/{domain}/{operation-id}/usecase.go)

## Input

EXAMPLE FOR POST

```json
{
    "username": "string" // required, min=5, max=50
    "email": "string" // required, email format
    "password": "string" // required, min=8
    "age": 12 // optional, integer
    "date_of_birth": "2001-01-31" // optional, date format
}
```

EXAMPLE FOR GET

- `username`: string, required, 3-50 chars
- `email`: string, required, email format
- `is_active`: bool, optional, default true

## Output

```json
{
    "id": "string",
    "username": "string",
    "email": "string",
    "is_active": bool,
    "age": 12, // nullable
    "date_of_birth": "2001-01-31" // nullable
}
```

## Execute

<!--
Describe WHAT the use case does, not HOW.
- Steps should map 1:1 to use case Execute() method logic
- Focus on business actions, not implementation details
- Do NOT include: logging, metrics, query mechanics, infrastructure concerns
-->

- Validate {what}

- Check {precondition}

- {Business action}

- Produce `{EventName}` event (if applicable)

- Return {result}
