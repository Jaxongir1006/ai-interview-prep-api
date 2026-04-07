package pguowbase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Jaxongir1006/ai-interview-prep-api/pkg/uowbase"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/observability/logger"
	"github.com/uptrace/bun"
)

// context key for lending/borrowing transactions between modules.
type ctxKey struct{}

func lend(ctx context.Context, idb bun.IDB) context.Context {
	return context.WithValue(ctx, ctxKey{}, idb)
}

func borrow(ctx context.Context) (bun.IDB, bool) {
	idb, ok := ctx.Value(ctxKey{}).(bun.IDB)
	return idb, ok
}

// Factory provides base methods for creating owned and borrowed UOW instances.
// Embed this in module-specific UOW factories.
type Factory struct {
	db         *bun.DB
	moduleName string
}

func NewFactory(db *bun.DB, moduleName string) Factory {
	return Factory{db: db, moduleName: moduleName}
}

// GenericFactory implements uowbase.Factory[T] using a constructor function.
// Modules only need to provide a constructor that wraps *Base into their concrete UOW type.
type GenericFactory[T uowbase.UnitOfWork] struct {
	base        Factory
	constructor func(*Base) T
}

// NewGenericFactory creates a factory that constructs module-specific UOW instances.
// The constructor receives a *Base and returns the module's concrete UOW type.
func NewGenericFactory[T uowbase.UnitOfWork](
	db *bun.DB,
	moduleName string,
	constructor func(*Base) T,
) *GenericFactory[T] {
	return &GenericFactory[T]{
		base:        NewFactory(db, moduleName),
		constructor: constructor,
	}
}

func (f *GenericFactory[T]) NewUOW(ctx context.Context) (T, error) {
	base, err := f.base.NewBase(ctx)
	if err != nil {
		var zero T
		return zero, err
	}
	return f.constructor(base), nil
}

func (f *GenericFactory[T]) NewBorrowed(ctx context.Context) (T, error) {
	base, err := f.base.NewBorrowedBase(ctx)
	if err != nil {
		var zero T
		return zero, err
	}
	return f.constructor(base), nil
}

// NewBase starts a new transaction and returns an owned Base.
func (f *Factory) NewBase(ctx context.Context) (*Base, error) {
	tx, err := f.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	return &Base{
		ctx:        ctx,
		idb:        tx,
		tx:         tx,
		owned:      true,
		moduleName: f.moduleName,
	}, nil
}

// NewBorrowedBase creates a Base from a transaction previously lent via Base.Lend().
func (f *Factory) NewBorrowedBase(ctx context.Context) (*Base, error) {
	idb, ok := borrow(ctx)
	if !ok {
		return nil, errx.New("no lent transaction found in context")
	}

	return &Base{
		ctx:        ctx,
		idb:        idb,
		owned:      false,
		moduleName: f.moduleName,
	}, nil
}

// Base provides the shared UOW functionality: transaction lifecycle and lending.
// Embed this in module-specific pgUOW structs and expose IDB() for repo construction.
type Base struct {
	ctx        context.Context
	idb        bun.IDB
	tx         bun.Tx
	owned      bool
	committed  bool
	moduleName string
}

// IDB returns the underlying bun.IDB for constructing repositories.
func (b *Base) IDB() bun.IDB {
	return b.idb
}

// Lend returns a context enriched with the underlying transaction,
// allowing other modules to borrow it via Factory.NewBorrowedBase().
func (b *Base) Lend() context.Context {
	return lend(b.ctx, b.idb)
}

// ApplyChanges commits the transaction. Returns an error if called on a borrowed UOW.
func (b *Base) ApplyChanges() error {
	if !b.owned {
		return errx.New("cannot commit a borrowed UOW — only the lender controls the transaction lifecycle")
	}
	err := errx.Wrap(b.tx.Commit())
	if err != nil {
		return errx.Wrap(err)
	}

	b.committed = true
	return nil
}

// DiscardUnapplied rolls back the transaction if not yet committed.
// Logs a warning if called on a borrowed UOW.
func (b *Base) DiscardUnapplied() {
	if !b.owned {
		logger.Named(b.moduleName+"_uow").
			WithContext(b.ctx).
			With("method", "DiscardUnapplied").
			Warn("cannot rollback a borrowed UOW — only the lender controls the transaction lifecycle")
		return
	}

	if b.committed {
		return
	}

	err := errx.Wrap(b.tx.Rollback())
	if err == nil || errors.Is(err, sql.ErrTxDone) {
		return
	}
	logger.Named(b.moduleName+"_uow").
		WithContext(b.ctx).
		With("method", "DiscardUnapplied").
		Warnx(err)
}
