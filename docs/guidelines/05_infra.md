# Infra

The infrastructure layer implements interfaces defined in the domain layer. It lives in `infra/` and is organized by technology: `infra/postgres/`, `infra/http/`, etc.

**Key principle:** interfaces are defined in the domain layer, infra just implements them. No business logic belongs here.

## PostgreSQL Repositories

We use the `repogen` package to build repositories. The goal is to keep repositories **minimal, simple, and obvious** — thin wrappers around repogen, not a place for business logic.

### Repogen-Only Repository

Most repositories need nothing beyond what repogen provides. No custom struct, no extra methods — just build and return:

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

### Domain Interface for Repogen-Only Repository

When a repository only uses repogen methods, the domain interface simply embeds `repogen.Repo`:

```go
type AdminRepo interface {
    repogen.Repo[Admin, AdminFilter]
}
```

### Custom Methods (Exceptional Cases)

Sometimes repogen's built-in methods are not enough — typically for **performance reasons** (e.g., bulk operations, specialized queries with joins). In such cases, add custom methods on top of repogen:

```go
// Domain interface — extends repogen with a custom method
type Repo interface {
    repogen.Repo[Session, Filter]

    DeleteExpired(ctx context.Context) (int64, error)
}
```

```go
// Infra implementation — embeds repogen, adds custom method
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

### When to Add Custom Methods

- **Performance** — repogen's general methods would be inefficient (e.g., `DeleteExpired` avoids loading all sessions just to delete them)
- **Complex queries** — joins, aggregations, or queries that don't fit repogen's filter model

### When NOT to Add Custom Methods

- **Filtering logic** — use repogen's Filter struct instead
- **Business logic** — keep it in use case or PBLC layer, not in the repository
- **Convenience** — don't create shortcut methods that just wrap repogen calls differently

### Filter Functions

Each repository provides a filter function that maps Filter struct fields to SQL conditions:

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

## Redis-Based Components

Components that use Redis (e.g., rate limiters, caches, distributed locks) follow the same pattern as PostgreSQL repositories: interface defined in the domain layer, implementation in `infra/redis/`.

## Kafka Publishers

Publishers (producers) follow the same domain/infra split. The publisher interface is defined in the domain layer, and the implementation lives in `infra/kafka/`. The domain container holds the publisher interface like any other dependency.

## HTTP Clients

HTTP clients are the exception to the domain/infra split. Because HTTP clients are often **shared between modules**, they live in the `pkg/` layer instead of inside a module's `infra/`:

```
pkg/clients/{client_name}/
├── client.go          # Interface and implementation
├── config.go          # Client configuration
├── fake.go            # Fake implementation for testing
└── client_test.go     # Unit tests
```

Since they already live in `pkg/`, their interface can be used directly without redefining it in the domain layer.
