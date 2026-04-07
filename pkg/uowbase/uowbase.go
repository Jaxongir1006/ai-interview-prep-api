package uowbase

import "context"

// UnitOfWork defines the shared lifecycle methods that every module's UnitOfWork interface should embed.
// Module-specific interfaces add their own repository accessors on top.
type UnitOfWork interface {
	// Lend returns a context enriched with the underlying transaction,
	// allowing other modules to borrow it via their Factory.NewBorrowed().
	Lend() context.Context

	// ApplyChanges commits the transaction.
	// Returns an error if called on a borrowed UOW.
	ApplyChanges() error

	// DiscardUnapplied rolls back the transaction if not yet committed.
	// Logs a warning if called on a borrowed UOW.
	DiscardUnapplied()
}

// Factory defines the shared interface for creating UOW instances.
// Modules define their own Factory by providing their concrete UnitOfWork type:
//
//	type Factory = uowbase.Factory[UnitOfWork]
type Factory[T UnitOfWork] interface {
	// NewUOW creates and returns a new UnitOfWork that owns its transaction.
	NewUOW(ctx context.Context) (T, error)

	// NewBorrowed creates a UnitOfWork that participates in an existing transaction
	// previously lent via UnitOfWork.Lend(). The borrowed UOW does not own the transaction.
	NewBorrowed(ctx context.Context) (T, error)
}
