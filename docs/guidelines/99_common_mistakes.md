# Common Mistakes

## Context Misusing

Every function that performs I/O or calls other services MUST accept `ctx context.Context` as the first parameter. Never create context inside a function.

```go
// Good
func FetchUser(ctx context.Context, id int64) (*User, error) {
    u, err := db.QueryUser(ctx, id)
    return u, errx.Wrap(err)
}

// Bad
func FetchUser(id int64) (*User, error) {
    ctx := context.Background()
    return db.QueryUser(ctx, id)
}
```

## Unnecessary Pointer Helpers

Use `&` for addressable values (struct fields, variables). Only use `lo.ToPtr()` for literals, constants, or function return values where `&` doesn't work.

```go
// Good — & works for struct fields and variables
repo.Get(ctx, Filter{ID: &in.ID})
repo.List(ctx, Filter{UserID: &userID})

// Bad — unnecessary dependency on lo
repo.Get(ctx, Filter{ID: lo.ToPtr(in.ID)})
repo.List(ctx, Filter{UserID: lo.ToPtr(userID)})

// Good — lo.ToPtr needed for non-addressable values
u.LastActiveAt = lo.ToPtr(time.Now())      // function return
filter.Order = lo.ToPtr(sorter.Asc)        // constant
```

## Unnecessary Pointer Slices

Slices are already reference types. Don't use pointer-to-struct slices (`[]*T`) unless you need to mutate elements through the slice. For return types and read-only data, use value slices.

```go
// Good — slice of values
func (uc *usecase) Execute(ctx context.Context, in *Request) ([]rbac.RolePermission, error)

// Bad — unnecessary pointer indirection
func (uc *usecase) Execute(ctx context.Context, in *Request) ([]*rbac.RolePermission, error)
```

## Temporal Comments

Never write comments that reference time, version history, or change context — e.g. `// Existing`, `// New`, `// Added in v2`, `// old implementation`, `// legacy`, `// was previously X`. These become meaningless after merge. Comments should describe **what** or **why**, not **when**.

```go
// Good
// Self-service
adminlogin.New(domainContainer)

// Good
// Enforces max 5 active sessions per user
err = uc.deleteExceededSessions(ctx, u, uow)

// Bad — temporal
// Existing
adminlogin.New(domainContainer)

// Bad — temporal
// New implementation replacing the old one
err = uc.deleteExceededSessions(ctx, u, uow)
```
