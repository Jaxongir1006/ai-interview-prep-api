# Use Cases

Everything about use cases: types and conventions.

## Use Case Types

| Type              | Trigger             | Controller          | Definition                              |
| ----------------- | ------------------- | ------------------- | --------------------------------------- |
| `UserAction`      | HTTP/gRPC request   | HTTP handler        | `ucdef.UserAction[*Request, *Response]` |
| `EventSubscriber` | Domain event        | Consumer            | `ucdef.EventSubscriber[*EventPayload]`  |
| `AsyncTask`       | Scheduler/on-demand | Taskmill worker     | `ucdef.AsyncTask[*Payload]`             |
| `ManualCommand`   | CLI                 | CLI handler (Cobra) | `ucdef.ManualCommand[*Input]`           |

### UserAction Example

```go
type Request struct {
    Username string `json:"username" validate:"required"`
    Password string `json:"password" validate:"required" mask:"true"`
}

type Response struct {
    AccessToken string `json:"access_token"`
}

type UseCase = ucdef.UserAction[*Request, *Response]

func (uc *usecase) OperationID() string { return "admin-login" }
```

### AsyncTask Example

```go
type Payload struct{}
type UseCase = ucdef.AsyncTask[*Payload]
func (uc *usecase) OperationID() string { return "clean-expired-sessions" }
```

## Use Case Conventions

- **Document first** — don't code before documenting
- **One package per UC** — `{module}/usecase/{domain}/{operation}/usecase.go`
- **OperationID** — every UC must implement `OperationID()` returning kebab-case name
- **Separate actors** — create separate UCs for different actor types (admin vs user)
- **No duplication** — if logic repeats, move to PBLC layer

## Implementation Must Reflect Documentation

The `Execute` method must mirror the documented execution steps. Each step from the documentation becomes a comment in the code, and the code block beneath it implements that step.

**Example** — `admin-login` use case:

```go
func (uc *usecase) Execute(ctx context.Context, in *Request) (*Response, error) {
    // Find admin by username
    a, err := uc.domainContainer.AdminRepo().Get(ctx, user.AdminFilter{...})
    ...

    // Check if admin is active
    if !a.IsActive {
        ...
    }

    // Verify password hash
    ok := hasher.Compare(in.Password, a.PasswordHash)
    ...

    // Start UOW
    uow, err := uc.domainContainer.UOWFactory().NewUOW(ctx)
    ...
    defer uow.DiscardUnapplied()

    // Enforce max active sessions limit
    err = uc.deleteExceededSessions(ctx, a, uow)
    ...

    // Create session record with tokens and meta info (IP, user_agent)
    s, err := uow.Session().Create(ctx, &session.Session{...})
    ...

    // Update admin's last_active_at timestamp
    ...

    // Apply UOW
    err = uow.ApplyChanges()
    ...

    // Return session tokens
    return &Response{...}, nil
}
```

Steps should describe **what** the system does, not **how** it does it. Focus on business intent, not implementation details.

- **Good:** `// Enforce max active sessions limit`
- **Bad:** `// Query sessions ordered by last_used_at ASC, calculate excess count, bulk delete oldest`

This keeps documentation and implementation in sync — reading the comments alone tells the full business flow.
