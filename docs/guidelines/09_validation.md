# Validation

Validation rules for each application layer.

## Validation by Layer

| Layer             | What to Validate                       | How                                                  |
| ----------------- | -------------------------------------- | ---------------------------------------------------- |
| Controller        | Input parameters only                  | `validate` tags on request structs (`validator/v10`) |
| Use Case          | Cross-field validation, business rules | Manual checks, repository lookups                    |
| PBLC              | All inputs strictly                    | Manual checks (doesn't know caller)                  |
| Repository (DB)   | Nothing                                | Validation is UC/PBLC responsibility                 |
| Repository (HTTP) | Inputs strictly                        | Catch errors before external calls                   |

## Controller Layer

For `user_action` types, input validation is automatic via `validate` tags:

```go
type Request struct {
    Username string `json:"username" validate:"required"`
    Password string `json:"password" validate:"required,min=8"`
}
```

## Use Case Layer

Handles validation not covered by controller:

- **Cross-field** — password confirmation matches
- **Business rules** — user has permission, object exists
- **Object-level** — verify IDs through repository

## PBLC Layer

- **Validate strictly** — PBLC doesn't know its caller
- **Return error codes** — not error types (UC assigns types)
- **Minimal codes** — only add new codes when truly necessary for callers

## Infrastructure Layer

- **Database repos** — no validation, trust UC/PBLC
- **External HTTP clients** — validate strictly before calling third-party services
