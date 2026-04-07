---
name: backend-guidelines
description: Use when implementing any backend feature, writing use cases, creating controllers, working with repositories (repogen), handling errors (errx), managing transactions (UOW), writing system tests, creating migrations (goose), fixing bugs, or making any code changes to the codebase.
---

# Backend Guidelines

## Overview

Reference guide for all backend development conventions in the Go enterprise blueprint project. Contains architecture rules, layer responsibilities, API design, coding standards, and implementation patterns. Every concept from the project's 16 guideline files is condensed here.

**Core principle:** Document first, implement layer-by-layer (bottom-up), test everything, never violate layer boundaries or cross-module import rules.

## Architecture

### File Structure

```
project-root/
├── cmd/              # Entry points (Cobra CLI)
├── config/           # YAML config per environment
├── internal/
│   ├── app/          # Bootstrap and lifecycle
│   ├── modules/      # Business modules
│   └── portal/       # Cross-module interfaces
├── pkg/              # Shared packages
├── migrations/       # DB migrations (goose)
└── tests/            # System and integration tests
```

### Module Structure

```
internal/modules/{module}/
├── module.go         # Initialization and wiring
├── domain/           # Entities, repo interfaces, container
│   └── {domain}/     # Grouped by business domain
├── usecase/          # One package per use case
│   └── {domain}/{operation}/usecase.go
├── pblc/             # Shared business logic components
├── infra/            # Implementations (postgres/, http/)
├── ctrl/             # Controllers (http/, cli/, consumer/, asynctask/)
└── embassy/          # Portal implementation
```

### Layers

| #   | Layer      | Location   | Responsibility                                 |
| --- | ---------- | ---------- | ---------------------------------------------- |
| 1   | Domain     | `domain/`  | Entities, value objects, repository interfaces |
| 2   | Use Case   | `usecase/` | Business operations, transaction boundaries    |
| 3   | PBLC       | `pblc/`    | Reusable business logic across use cases       |
| 4   | Infra      | `infra/`   | Repository implementations, external clients   |
| 5   | Controller | `ctrl/`    | HTTP handlers, CLI, consumers, async tasks     |

**Rules:**

- Business logic in use case or PBLC only — never in controllers or repos
- Each layer depends only on layers above (controller -> UC -> PBLC -> domain <- infra)
- Infrastructure implements domain interfaces

### Cross-Module Communication

Modules communicate **only** through Portal interfaces. No direct imports between modules.

- **Portal** (`internal/portal/{module}/`) — interface contract
- **Embassy** (`internal/modules/{module}/embassy/`) — implementation
- Each module owns its data — no cross-schema joins
- Embassies wired in `internal/app/` during startup

## API Design

HTTP API using only **GET** (queries) and **POST** (mutations). Not REST.

**URL:** `{method} api/v1/{module}/{operation-id}`

**Rules:**

- No path parameters — use query params (GET) or JSON body (POST)
- Operation ID = use case name in kebab-case
- GET inputs via query parameters only
- POST inputs via JSON body only (exception: file uploads)

**Response formats:**

- List: `{ "content": [] }` — always wrap, never bare arrays
- Paginated: `{ "page_number": 1, "page_size": 20, "count": 150, "content": [] }`
- Error: `{ "trace_id": "...", "error": { "code": "...", "message": "...", "cause": "...", "fields": {}, "details": {} } }`
- Header: `X-Trace-ID` on every response

## Code Style

**Style guides:** Effective Go, Go Code Review Comments, Uber Go Style Guide

**Naming:**

- Packages: lowercase, single-word, matches directory (`user`, `adminlogin`)
- Files: snake_case (`user_repository.go`)
- Operation IDs: kebab-case (`admin-login`)

**Structs:** Exported fields first, group related fields, don't repeat default bun tags

**Comments:**

- Start exported item comments with the item name
- No obvious comments (`// GetUser gets a user`)
- No action/narrative comments (`// I fixed this`, `// Added this to handle...`)
- No temporal comments (`// Existing`, `// New`, `// Added in v2`, `// legacy`)

**Error handling style** — always separate call from check:

```go
// Good
err := doSomething()
if err != nil {
    return errx.Wrap(err)
}

// Bad — no inline if-assignments
if err := doSomething(); err != nil {
    return errx.Wrap(err)
}
```

**Formatting:** golangci-lint (config at `.golangci.yml`)

## Development Workflow

**Document-first** at every level:

| Level    | Document first             | Then implement             |
| -------- | -------------------------- | -------------------------- |
| System   | Design flows               | Break into modules and UCs |
| Module   | Write overview + ERD       | Break into use cases       |
| Use case | Write UC doc from template | Implement                  |

**The cycle:** Analyze -> Document -> Implement -> Test -> Review & Verify

**Implementation order** (bottom-up):

1. Migrations -> 2. Domain -> 3. Infra -> 4. PBLC -> 5. Use Case -> 6. Controller -> 7. DI Containers

**Feedback loop:** If implementation reveals missing edge cases, update docs BEFORE coding the fix. Never implement undocumented behavior.

### Review & Verify

**Sync checks:**

- Docs <-> Code: execute steps match `Execute` method comments
- Docs <-> Tests: every documented behavior has a test
- Docs <-> Docs: UC index in README.md is up-to-date
- Lint: `make lint` passes
- Tests: `make test` and `make test-system` pass

**Quality checks:**

- No layer violations, no cross-module imports without portals
- No SQL injection, no sensitive data in logs
- No unnecessary abstractions, no N+1 queries
- UOW not held open during read-only operations

## Documentation

- **Source of truth** — don't code before documenting
- **UC docs = API specs** — frontend devs should need no clarification
- **Templates** in `docs/specs/templates/usecases/` (user_action, async_task, event_subscriber, manual_command)
- **Input/Output format** — use JSON code blocks with inline `//` comments for validation rules, NOT markdown tables:
  ```json
  {
    "username": "string", // required, min=3, max=50
    "password": "string" // required, min=8
  }
  ```
- **Output fields** — document EVERY field the API returns, including timestamps and FKs
- **Validation rules** — specify types, min/max, allowed values, required vs optional, nullable (`// nullable`)
- **Sort fields** — when a UC exposes dynamic sort, list allowed sortable fields and default: `sort: string, optional — sortable fields: name, created_at. Default: created_at:desc`
- **Transaction boundaries** — mark with "Start UOW" and "Apply UOW" steps. Place Start after read-only checks, Apply after last write.
- **UC Index** — register every UC in README.md table

### Documentation Structure

```
docs/specs/modules/{module}/
├── overview.md      # Purpose, responsibilities, main entities
├── ERD.md           # Mermaid ERD with column types
└── usecases/{domain}/{operation}.md
```

ERD rules: follow normalization, describe non-obvious fields, skip universal fields (id, created_at, updated_at). Use `timestamptz` for all timestamps.

## Use Cases

### Types

| Type              | Trigger             | Controller    | Definition                              |
| ----------------- | ------------------- | ------------- | --------------------------------------- |
| `UserAction`      | HTTP/gRPC           | HTTP handler  | `ucdef.UserAction[*Request, *Response]` |
| `EventSubscriber` | Domain event        | Consumer      | `ucdef.EventSubscriber[*EventPayload]`  |
| `AsyncTask`       | Scheduler/on-demand | Taskmill      | `ucdef.AsyncTask[*Payload]`             |
| `ManualCommand`   | CLI                 | Cobra handler | `ucdef.ManualCommand[*Input]`           |

### Conventions

- **One package per UC:** `{module}/usecase/{domain}/{operation}/usecase.go`
- **OperationID:** every UC implements `OperationID()` returning kebab-case name
- **Document first** — don't code before documenting
- **Separate actors** — create separate UCs for different actor types (admin vs user)
- **No duplication** — repeated logic moves to PBLC layer

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

### Execute Must Mirror Documentation

Each documented step becomes a comment in code. The code block beneath implements that step.

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

    // Create session record with tokens and meta info
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

Steps describe **what** not **how**:

- Good: `// Enforce max active sessions limit`
- Bad: `// Query sessions ordered by last_used_at ASC, calculate excess count, bulk delete oldest`

## Controllers

### Rules

- One-to-one mapping with use cases
- No business logic — delegate everything to use cases
- Always use forwarders (`forward.To*`, `worker.ForwardTo*`)
- Manual handlers only for file uploads/downloads

### Wiring Reference

| Type      | Wiring                                                                                           |
| --------- | ------------------------------------------------------------------------------------------------ |
| HTTP      | `v1.Post("/path", forward.ToUserAction(c.usecaseContainer.SomeUseCase()))`                       |
| Consumer  | `kafka.NewConsumer(brokerCfg, cfg, forward.ToEventSubscriber(c.usecaseContainer.SomeUseCase()))` |
| AsyncTask | `worker.ForwardToAsyncTask(c.worker, c.usecaseContainer.SomeUseCase())`                          |
| CLI       | `c.usecaseContainer.SomeUseCase().Execute(ctx, input)`                                           |

## PBLC (Packaged Business Logic Components)

Use cases are command-based separate structs — no shared struct for common logic. PBLC provides the shared layer.

### When to Use

- **Deduplication** — logic repeating across multiple use cases
- **Complex business logic** — state machines, strategy patterns

### Design Freedom

No prescribed structure. Choose what fits:

- **Simple functions** — for deduplication, accept dependencies as params
- **Structs with interfaces** — for complex encapsulated components
- **OOP patterns** — when business logic genuinely calls for it

### Rules

- Called only from use cases — never from controllers
- Validate all inputs strictly — PBLC doesn't know its caller
- Return error codes — use case layer assigns error types

## Infrastructure

### PostgreSQL Repositories

Use `repogen` package. Keep repos minimal — thin wrappers, not business logic.

**Repogen-only repo** (most common — no custom struct needed):

```go
func NewAdminRepo(idb bun.IDB) user.AdminRepo {
    return repogen.NewPgRepoBuilder[user.Admin, user.AdminFilter](idb).
        WithSchemaName(schemaName).
        WithNotFoundCode(user.CodeAdminNotFound).
        WithConflictCodesMap(map[string]string{
            "admins_username_key": user.CodeAdminUsernameConflict,
        }).
        WithFilterFunc(adminFilterFunc).
        Build()
}
```

**Domain interface for repogen-only:**

```go
type AdminRepo interface {
    repogen.Repo[Admin, AdminFilter]
}
```

**Custom methods** — only for performance (bulk ops, specialized queries):

```go
// Domain — extend repogen
type Repo interface {
    repogen.Repo[Session, Filter]
    DeleteExpired(ctx context.Context) (int64, error)
}

// Infra — embed repogen, add custom
type sessionRepo struct {
    repogen.Repo[session.Session, session.Filter]
    idb bun.IDB
}

func (r *sessionRepo) DeleteExpired(ctx context.Context) (int64, error) {
    res, err := r.idb.NewDelete().
        Model((*session.Session)(nil)).
        Where("refresh_token_expires_at < ?", time.Now()).
        Exec(ctx)
    ...
}
```

**When to add custom methods:** Performance, complex queries (joins, aggregations)
**When NOT to:** Filtering logic, business logic, convenience wrappers

### Filter Functions

Map Filter struct fields to SQL conditions:

```go
func adminFilterFunc(q *bun.SelectQuery, f user.AdminFilter) *bun.SelectQuery {
    if f.ID != nil {
        q = q.Where("id = ?", *f.ID)
    }
    if f.Username != nil {
        q = q.Where("username = ?", *f.Username)
    }
    ...
    return q
}
```

### Redis / Kafka

Same domain/infra split: interface in domain, implementation in `infra/redis/` or `infra/kafka/`. Domain container holds the interface.

### HTTP Clients

Exception — live in `pkg/` (shared between modules):

```
pkg/clients/{client_name}/
├── client.go      # Interface + implementation
├── config.go      # Configuration
├── fake.go        # Fake for testing
└── client_test.go # Unit tests
```

Use directly without redefining in domain layer.

## DI Containers

Each layer (except controllers) provides a container.

| Layer    | Holds                                      |
| -------- | ------------------------------------------ |
| Domain   | Repositories, UOW factory, pkg/ interfaces |
| PBLC     | PBLC component instances                   |
| Use Case | Use case instances                         |

Controllers don't provide containers — they are entry points, not dependencies.

### Pattern

All containers follow: unexported fields, `NewContainer(...)` constructor, getter methods.

```go
type Container struct {
    adminRepo   user.AdminRepo
    sessionRepo session.Repo
    uowFactory  uow.Factory
}

func NewContainer(
    adminRepo user.AdminRepo,
    sessionRepo session.Repo,
    uowFactory uow.Factory,
) *Container {
    return &Container{adminRepo, sessionRepo, uowFactory}
}

func (c *Container) AdminRepo() user.AdminRepo { return c.adminRepo }
func (c *Container) SessionRepo() session.Repo { return c.sessionRepo }
func (c *Container) UOWFactory() uow.Factory   { return c.uowFactory }
```

Shared `pkg/` interfaces go directly in domain container — don't redefine.

## Error Handling

Use `github.com/code19m/errx` for ALL error handling.

### Core Rules

- Always handle errors, never ignore
- Use `errx.Wrap()` to preserve stack traces
- Use error codes for programmatic handling
- Never `panic`

### Layer Responsibilities

| Layer      | Error Types                        | Error Codes                         |
| ---------- | ---------------------------------- | ----------------------------------- |
| Controller | Handled by framework automatically | —                                   |
| Use Case   | Assigns types (knows the actor)    | Checks from downstream, defines own |
| PBLC       | Never sets (default internal)      | May return when callers need them   |
| Infra      | Never sets (default internal)      | May return (not found, conflict)    |

**Only the use case layer assigns error types.** Downstream layers always return `errx.T_Internal` (default).

### Error Codes

Define as constants in domain layer:

```go
const (
    CodeAdminNotFound  = "ADMIN_NOT_FOUND"
    CodeIncorrectCreds = "INCORRECT_CREDENTIALS"
)
```

**Check codes:**

```go
if errx.IsCodeIn(err, user.CodeAdminNotFound) {
    // branch on specific error
}
```

**Reassign type by code:**

```go
return errx.WrapWithTypeOnCodes(err, errx.T_NotFound, user.CodeAdminNotFound)
```

Use `WrapWithTypeOnCodes` when repo call uses **user-provided input** (caller supplied an ID that doesn't exist -> not found for caller). For **internal data** (ID from another query), just `errx.Wrap(err)` — it's an internal consistency error.

### Inline Return

All errx wrap functions return `nil` on `nil` input. When wrapping is the last operation:

```go
// Good — clean, final return
role, err := repo.Create(ctx, &rbac.Role{Name: in.Name})
return role, errx.WrapWithTypeOnCodes(err, errx.T_Conflict, rbac.CodeRoleNameConflict)

// Bad — unnecessary nil check at final return
role, err := repo.Create(ctx, &rbac.Role{Name: in.Name})
if err != nil {
    return nil, errx.WrapWithTypeOnCodes(err, errx.T_Conflict, rbac.CodeRoleNameConflict)
}
return role, nil
```

Only for **final return**. Mid-function errors still need `if err != nil`.

### Creating Errors

```go
return nil, errx.New(
    "admin not found",
    errx.WithType(errx.T_NotFound),
    errx.WithCode(user.CodeAdminNotFound),
)
```

Pre-defined error vars only when same error reused multiple times within a UC:

```go
var errIncorrectCreds = errx.New(
    "username or password is incorrect",
    errx.WithType(errx.T_Validation),
    errx.WithCode(user.CodeIncorrectCreds),
)

// Usage with context to distinguish cause:
return nil, errx.Wrap(errIncorrectCreds, errx.WithDetails(errx.D{"cause": "password"}))
```

### Error Details

Best added when creating the error. When adding structs to details, **mask first** using the `mask` package — sensitive fields tagged `mask:"true"` will be hidden.

## Validation

| Layer             | What                        | How                                 |
| ----------------- | --------------------------- | ----------------------------------- |
| Controller        | Input parameters            | `validate` tags (`validator/v10`)   |
| Use Case          | Cross-field, business rules | Manual checks, repo lookups         |
| PBLC              | All inputs strictly         | Manual checks (doesn't know caller) |
| DB Repos          | Nothing                     | Trust UC/PBLC                       |
| HTTP Client Repos | Inputs strictly             | Catch errors before external calls  |

Controller validation example:

```go
type Request struct {
    Username string `json:"username" validate:"required"`
    Password string `json:"password" validate:"required,min=8"`
}
```

PBLC returns error codes, not types — UC assigns types.

## Transaction Management (UOW)

Unit of Work pattern. Transactions at **use case/PBLC layer only**, never at repo/infra.

**Rule: who opens a transaction must close it.**

### How It Works

Repos accept `bun.IDB` — works with both `*bun.DB` (direct) and `bun.Tx` (transaction). Same constructor for both. The repo doesn't know which one it gets.

### Shared Packages

- `pkg/uowbase` — domain interfaces (`UnitOfWork`, `Factory[T]`)
- `pkg/uowbase/pguowbase` — PostgreSQL implementation (`Base`, `Factory`)

### Domain Interface

Each module defines its own UOW by embedding `uowbase.UnitOfWork` and adding repo accessors:

```go
type Factory = uowbase.Factory[UnitOfWork]

type UnitOfWork interface {
    uowbase.UnitOfWork
    Admin() user.AdminRepo
    Session() session.Repo
}
```

### Infra Implementation

```go
func NewUOWFactory(db *bun.DB) uow.Factory {
    return pguowbase.NewGenericFactory(db, schemaName,
        func(base *pguowbase.Base) uow.UnitOfWork {
            return &pgUOW{Base: base}
        },
    )
}

type pgUOW struct{ *pguowbase.Base }

func (u *pgUOW) Admin() user.AdminRepo   { return NewAdminRepo(u.IDB()) }
func (u *pgUOW) Session() session.Repo   { return NewSessionRepo(u.IDB()) }
```

### Usage in Use Cases

```go
// Start UOW
uow, err := uc.domainContainer.UOWFactory().NewUOW(ctx)
if err != nil {
    return errx.Wrap(err)
}
defer uow.DiscardUnapplied()

// Use transactional repos
session, err := uow.Session().Create(ctx, &session.Session{...})
if err != nil {
    return errx.Wrap(err)
}

// Commit
err = uow.ApplyChanges()
return errx.Wrap(err)
```

### Critical Mistake: Mixing Repos

Once UOW is started, **always use UOW repos**. Never mix domain container repos with UOW repos:

```go
// WRONG — admin update outside transaction
_, err = uow.Session().Create(ctx, &session)          // inside tx
_, err = uc.domainContainer.AdminRepo().Update(ctx, a) // outside tx!

// CORRECT — all through UOW
_, err = uow.Session().Create(ctx, &session)  // inside tx
_, err = uow.Admin().Update(ctx, a)           // inside tx
```

### Lend / Borrow (Cross-Module Transactions)

When a portal call needs the same transaction:

**Lender (owns transaction):**

```go
uow, err := uc.domainContainer.UOWFactory().NewUOW(ctx)
defer uow.DiscardUnapplied()

// ... own work ...

// Lend transaction before portal call
err = uc.portalContainer.Audit().RecordAction(uow.Lend(), ...)

err = uow.ApplyChanges()
```

**Borrower (embassy):**

```go
func (e *embassy) RecordAction(ctx context.Context, ...) error {
    uow, err := e.uowFactory.NewBorrowed(ctx)
    // NO defer DiscardUnapplied — borrowed must not rollback
    // NO ApplyChanges — borrowed must not commit

    // ... write using uow repos — within lender's transaction ...

    return nil
}
```

### UOW Rules

- Always `defer uow.DiscardUnapplied()` after creating owned UOW
- Use UOW repos for transactional ops, domain container repos for read-only
- Only use UOW when multiple writes need atomicity — single writes don't need it
- Never `ApplyChanges()` or `DiscardUnapplied()` on borrowed UOW
- Use `Lend()` / `NewBorrowed()` for cross-module transaction sharing through portals

## List Manipulations

### Response Structure

All lists wrapped in `{ "content": [] }` — even non-paginated. Never bare arrays.

Prefer returning full domain entities over cherry-picked response DTOs. Create dedicated structs only when joining data or computing fields.

### Pagination

Use `rise-and-shine/pkg/pagination`.

```go
type Request struct {
    pagination.Request
    Name   *string `query:"name"`
    Status *string `query:"status"`
}
```

```go
func (uc *usecase) Execute(ctx context.Context, in *Request) (*pagination.Response[SomeDTO], error) {
    // Normalize pagination params
    in.Normalize()  // defaults: page_number=1, page_size=20, max 100

    filter := entity.Filter{
        Name:   in.Name,
        Limit:  lo.ToPtr(in.PageSize()),
        Offset: lo.ToPtr(in.Offset()),
    }

    items, count, err := uc.domainContainer.SomeRepo().ListWithCount(ctx, filter)
    if err != nil {
        return nil, errx.Wrap(err)
    }

    return pagination.NewResponse(items, count, in.Request), nil
}
```

### Sorting

**Fixed-field** (internal, not API-exposed):

```go
// Domain Filter
type Filter struct {
    OrderByLastUsedAt *sorter.SortDirection
}

// Use case
sessions, err := repo.List(ctx, session.Filter{
    OrderByLastUsedAt: lo.ToPtr(sorter.Asc),
})

// Infra filter function
if f.OrderByLastUsedAt != nil {
    q = q.Order("last_used_at " + cast.ToString(f.OrderByLastUsedAt))
}
```

**Dynamic** (user-facing, `rise-and-shine/pkg/sorter`):

```go
// Request
Sort string `query:"sort"`  // e.g. "name:asc,created_at:desc"

// Use case — parse with allowlist
SortOpts: sorter.MakeFromStr(in.Sort, "name", "created_at", "status")

// Infra filter function
for _, o := range f.SortOpts {
    q = q.Order(o.ToSQL())
}
```

`MakeFromStr` parses, validates against allowlist, silently drops unknowns. Returns `nil` for empty string.

**Rules:** Every Filter has `SortOpts sorter.SortOpts`. Fixed-field checked before SortOpts in filter function. Always pass an allowlist to `MakeFromStr`. When a UC exposes dynamic sort, the UC doc **must** list allowed sortable fields and default sort.

### Filter Struct Convention

```go
type Filter struct {
    // Exact-match (pointer = optional, nil = don't filter)
    ID       *string
    Username *string
    IsActive *bool

    // Multi-value (nil = don't filter, empty = match nothing)
    IDs []int64

    // Pagination (always *int)
    Limit  *int
    Offset *int

    // Fixed-field sorting (optional)
    OrderByCreatedAt *sorter.SortDirection

    // Dynamic sorting (always present)
    SortOpts sorter.SortOpts
}
```

Multi-value: guard with `if f.IDs != nil` not `if len(f.IDs) > 0`.
Filter structs belong in domain layer. Filter functions belong in infra layer.

### Filter Function Order

```go
func filterFunc(q *bun.SelectQuery, f entity.Filter) *bun.SelectQuery {
    // 1. WHERE conditions
    if f.ID != nil {
        q = q.Where("id = ?", *f.ID)
    }
    if f.IDs != nil {
        q = q.Where("id IN (?)", bun.In(f.IDs))
    }

    // 2. Pagination
    if f.Limit != nil {
        q = q.Limit(*f.Limit)
    }
    if f.Offset != nil {
        q = q.Offset(*f.Offset)
    }

    // 3. Fixed-field sorting
    if f.OrderByCreatedAt != nil {
        q = q.Order("created_at " + cast.ToString(f.OrderByCreatedAt))
    }

    // 4. Dynamic sorting (always last)
    for _, o := range f.SortOpts {
        q = q.Order(o.ToSQL())
    }

    return q
}
```

### Search

`Search` is a `string` (not `*string`), empty = no search:

```go
// Request and Filter
Search string `query:"search"`

// Filter function
if f.Search != "" {
    q = q.WhereGroup(" AND ", func(sq *bun.SelectQuery) *bun.SelectQuery {
        return sq.
            Where("name ILIKE ?", "%"+f.Search+"%").
            WhereOr("description ILIKE ?", "%"+f.Search+"%")
    })
}
```

Always wrap with `%` in filter function. Use `WhereGroup` with OR across columns. Search and exact-match filters combine with AND.

### Joining Rules

1. **Prefer UC-layer merging** — query each repo separately, merge in use case
2. **Within-module joins** — allowed as last resort (filter/sort by joined column, aggregations, performance)
3. **Cross-module joins** — NEVER. Use portals to get data from other modules.

Cross-module pattern:

```go
products, _ := uc.domainContainer.ProductRepo().List(ctx, filter)

creatorIDs := uniqueIDs(products, func(p Product) string { return p.CreatedBy })
creators, _ := uc.portalContainer.Auth().GetAdminsByIDs(ctx, creatorIDs)

return mergeProductsWithCreators(products, creators), nil
```

## Testing

### Unit Tests — `pkg/` Layer

Required for HTTP clients, shared utilities, reusable components. Use `net/http/httptest` for HTTP clients. Test files use `_test.go` suffix with `_test` package name.

### System Tests

All use cases must have **100% system test coverage**.

#### Directory Structure

```
tests/
├── state/
│   ├── database/             # DB helpers (GetTestDB, Empty)
│   └── {module}/             # Module state helpers (Given*, Get*)
└── system/
    ├── trigger/              # UC trigger helpers
    └── modules/{module}/{domain}/{operation}_test.go
```

#### GIVEN-WHEN-THEN

```go
func TestUseCaseName(t *testing.T) {
    // GIVEN
    database.Empty(t)
    admins := auth.GivenAdmins(t, map[string]any{})

    // WHEN
    resp := trigger.UserAction(t).POST("/api/v1/endpoint").
        WithJSON(payload).Expect()

    // THEN
    resp.Status(http.StatusOK)
    assert.Equal(t, expected, actual)
}
```

#### Trigger Functions

| UC Type            | Trigger                                      |
| ------------------ | -------------------------------------------- |
| `user_action`      | `trigger.UserAction(t).POST(...)`            |
| `manual_command`   | `trigger.ManualCommand(t, args...)`          |
| `async_task`       | `trigger.AsyncTask(t, queue, opID, payload)` |
| `event_subscriber` | `trigger.EventSubscriber(t, topic, event)`   |

#### State Helpers

**Given** — create test data:

```go
auth.GivenAdmins(t, map[string]any{})                                                     // defaults
auth.GivenAdmins(t, map[string]any{"username": "custom", "is_active": false})              // custom
auth.GivenAdmins(t, map[string]any{"username": "alice"}, map[string]any{"username": "bob"}) // multiple
```

**Getters** — verify state:

```go
auth.GetAdminByUsername(t, "alice")
auth.AdminExists(t, "alice")
auth.SessionCount(t, adminID)
auth.HasPermission(t, "admin", id, "perm")
```

**Passwords:** Use `auth.TestPassword1` — pre-computed hash to avoid bcrypt overhead.

#### Test Rules

- **Isolation** — `database.Empty(t)` at start of each test
- **Independence** — no dependency on other tests' data
- **Deterministic** — same results every run
- **Table-driven** — always prefer where possible
- **Comprehensive success** — one test covers the full success scenario (response, DB state, side effects)
- **Minimize count** — each test should be meaningful and verify as much as possible

#### Deriving Tests from Documentation

| Doc Element           | Test Case                                                   |
| --------------------- | ----------------------------------------------------------- |
| Success execute steps | One comprehensive success test                              |
| Simple failures       | One table-driven test (validation, missing fields)          |
| Complex failures      | Dedicated test per scenario (inactive user, limit exceeded) |
| Start/Apply UOW       | Verify atomicity (partial failure = no state change)        |

#### Commands

```bash
make test        # Unit tests (pkg/)
make test-system # System tests
```

## Common Mistakes

- **Context misusing** — every I/O function MUST accept `ctx context.Context` as first param. Never create context inside functions.
- **Unnecessary pointer helpers** — use `&` for addressable values. Use `lo.ToPtr()` only for literals, constants, or function returns where `&` doesn't work.
- **Unnecessary pointer slices** — use `[]T` not `[]*T` for read-only data. Slices are already reference types.
- **Temporal comments** — never reference time/version (`// old`, `// legacy`, `// was previously X`).

## Observability

- **Logger:** `github.com/rise-and-shine/pkg/observability/logger` — no other logger packages
- **Framework handles:** HTTP logging, error alerting, tracing (all automatic)
- **Manual logging for:** HTTP clients (debug), business events in UCs, app lifecycle
- **API:** `logger.WithContext(ctx).Named("auth_usecase").With("key", val).Info("msg")`
- **Levels:** debug (HTTP client detail), info (operational), warn, error
- **Tracing:** OpenTelemetry, auto-initialized. **Metrics:** OpenTelemetry Metrics.
- **Alerting:** `github.com/rise-and-shine/pkg/observability/alert`, provider-based (Telegram)

## DB Migrations

- **Package:** goose
- **Commands:** `make migrate-create`, `make migrate-up`, `make migrate-down`
- **Naming:** module prefix, snake_case (`auth_init_schema`, `platform_init_taskmill`)
- **Single folder:** `./migrations`, no subfolders
- **Query order:** CREATE TABLE -> CREATE INDEX -> ALTER TABLE (foreign keys)
- **Rollback:** one at a time, reverse order, use `IF EXISTS`
- **Timestamps:** always `timestamptz`
- **Auto-execution** on app startup
