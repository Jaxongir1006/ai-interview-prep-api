# Transaction Management

Transaction management using the Unit of Work (UOW) pattern.

## Core Principle

Transactions are managed at the **use case or PBLC layer**, never at the repository or infrastructure layer. Repositories are unaware of whether they operate on a direct DB connection or a transaction.

**Rule: who opens a transaction must close it.**

## How It Works

Repositories accept `bun.IDB` — an interface implemented by both `*bun.DB` (direct connection) and `bun.Tx` (transaction). This means the same repository constructor works in both cases:

- **Without UOW** — domain container creates repos with the DB connection: `NewAdminRepo(db)`
- **With UOW** — UOW creates repos with the transaction: `NewAdminRepo(tx)`

The repository doesn't know and doesn't care which one it gets.

## Shared Packages

Transaction management is split into shared packages to avoid duplication across modules:

- **`pkg/uowbase`** — domain-level interfaces (`UnitOfWork`, `Factory[T]`) that module-specific interfaces embed
- **`pkg/uowbase/pguowbase`** — PostgreSQL/bun implementation (`Base`, `Factory`) that module-specific infra embeds

## Domain Interface

Each module defines its own `UnitOfWork` by embedding `uowbase.UnitOfWork` and adding repository accessors. The `Factory` is a type alias over the generic `uowbase.Factory[T]`:

```go
// Factory — type alias with the module's concrete UnitOfWork type
type Factory = uowbase.Factory[UnitOfWork]

// UnitOfWork — embeds shared lifecycle, adds module-specific repos
type UnitOfWork interface {
    uowbase.UnitOfWork

    Admin() user.AdminRepo
    Session() session.Repo
    ...
}
```

## Infra Implementation

Use `pguowbase.NewGenericFactory` to create the factory — it handles `NewUOW` and `NewBorrowed` via a constructor function. The module only defines `pgUOW` with repository accessors:

```go
func NewUOWFactory(db *bun.DB) uow.Factory {
	return pguowbase.NewGenericFactory(
		db,
		schemaName,
		func(base *pguowbase.Base) uow.UnitOfWork {
			return &pgUOW{Base: base}
		},
	)
}

type pgUOW struct {
    *pguowbase.Base
}

func (u *pgUOW) Admin() user.AdminRepo {
    return NewAdminRepo(u.IDB())
}
```

## Usage in Use Cases

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

## Common Mistake: Mixing UOW and Domain Container Repos

Once a UOW is started, **always use UOW repos** for the rest of that operation. Never mix domain container repos (direct DB connection) with UOW repos (transaction) — the domain container repos operate outside the transaction and their changes won't be part of the atomic commit.

```go
// WRONG — admin update happens outside the transaction
uow, err := uc.domainContainer.UOWFactory().NewUOW(ctx)
...
_, err = uow.Session().Create(ctx, &session)          // inside transaction
_, err = uc.domainContainer.AdminRepo().Update(ctx, a) // outside transaction!

// CORRECT — all operations go through UOW
uow, err := uc.domainContainer.UOWFactory().NewUOW(ctx)
...
_, err = uow.Session().Create(ctx, &session)  // inside transaction
_, err = uow.Admin().Update(ctx, a)           // inside transaction
```

## Cross-Module Transaction Sharing (Lend / Borrow)

When a use case calls another module through a portal and both modules need to participate in the same transaction, use the **Lend/Borrow** pattern.

**Lend** — the originating UOW puts its transaction into the context. The UOW stores the context from creation time, so `Lend()` takes no parameters.

**NewBorrowed** — the receiving module's UOW Factory creates a UOW from the lent transaction. The borrowed UOW does not own the transaction — calling `ApplyChanges()` returns an error and calling `DiscardUnapplied()` logs a warning. Only the originating module controls the transaction lifecycle.

```go
// Module A's use case — owns the transaction
uow, err := uc.domainContainer.UOWFactory().NewUOW(ctx)
if err != nil {
    return errx.Wrap(err)
}
defer uow.DiscardUnapplied()

// ... Module A's own work ...

// Lend transaction before cross-module portal call
err = uc.portalContainer.Audit().RecordAction(uow.Lend(), ...)
if err != nil {
    return errx.Wrap(err)
}

err = uow.ApplyChanges()
return errx.Wrap(err)
```

```go
// Module B's embassy (portal implementation) — borrows the transaction
func (e *embassy) RecordAction(ctx context.Context, ...) error {
    uow, err := e.uowFactory.NewBorrowed(ctx)
    if err != nil {
        return errx.Wrap(err)
    }
    // NO defer uow.DiscardUnapplied() — borrowed UOW must not rollback
    // NO uow.ApplyChanges() — borrowed UOW must not commit

    // ... write using uow repos — all within Module A's transaction ...

    return nil
}
```

**Rule: who opens a transaction must close it.** Borrowed UOWs never commit or rollback — that's the lender's responsibility.

## Rules

- Always `defer uow.DiscardUnapplied()` immediately after creating an owned UOW
- Use `uow.Repo()` methods for transactional operations
- Call `uow.ApplyChanges()` to commit
- Errors before `ApplyChanges()` trigger automatic rollback via the deferred `DiscardUnapplied()`
- Use UOW only when multiple writes need to be atomic — for single read operations, use domain container repos directly
- Use `Lend()` / `NewBorrowed()` for cross-module transaction sharing through portals
- Never call `ApplyChanges()` or `DiscardUnapplied()` on a borrowed UOW
