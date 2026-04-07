package pguowbase_test

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"

	"go-enterprise-blueprint/pkg/uowbase"
	"go-enterprise-blueprint/pkg/uowbase/pguowbase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

// fakeDriver is a minimal SQL driver that supports transactions for testing.
type fakeDriver struct{}

func (d *fakeDriver) Open(_ string) (driver.Conn, error) {
	return &fakeConn{}, nil
}

type fakeConn struct{}

var errNotImplemented = errors.New("not implemented")

func (c *fakeConn) Prepare(_ string) (driver.Stmt, error) { return nil, errNotImplemented }
func (c *fakeConn) Close() error                          { return nil }
func (c *fakeConn) Begin() (driver.Tx, error)             { return &fakeTx{}, nil }

type fakeTx struct{ committed, rolledBack bool }

func (t *fakeTx) Commit() error   { t.committed = true; return nil }
func (t *fakeTx) Rollback() error { t.rolledBack = true; return nil }

func init() { //nolint:gochecknoinits // allow here for testing
	sql.Register("pguowbase_test", &fakeDriver{})
}

func newTestDB(t *testing.T) *bun.DB {
	t.Helper()
	sqlDB, err := sql.Open("pguowbase_test", "")
	require.NoError(t, err)
	return bun.NewDB(sqlDB, pgdialect.New())
}

func TestNewBase(t *testing.T) {
	t.Run("creates an owned base", func(t *testing.T) {
		factory := pguowbase.NewFactory(newTestDB(t), "test")

		base, err := factory.NewBase(t.Context())
		require.NoError(t, err)
		assert.NotNil(t, base.IDB())
	})
}

func TestNewBorrowedBase(t *testing.T) {
	t.Run("returns error when no transaction was lent", func(t *testing.T) {
		factory := pguowbase.NewFactory(newTestDB(t), "test")

		base, err := factory.NewBorrowedBase(t.Context())
		require.Error(t, err)
		assert.Nil(t, base)
		assert.Contains(t, err.Error(), "no lent transaction found in context")
	})

	t.Run("creates a borrowed base from lent context", func(t *testing.T) {
		factory := pguowbase.NewFactory(newTestDB(t), "test")

		owned, err := factory.NewBase(t.Context())
		require.NoError(t, err)

		lentCtx := owned.Lend()

		borrowed, err := factory.NewBorrowedBase(lentCtx)
		require.NoError(t, err)
		assert.NotNil(t, borrowed.IDB())
	})
}

func TestLend(t *testing.T) {
	t.Run("lent context allows borrowing the same IDB", func(t *testing.T) {
		factory := pguowbase.NewFactory(newTestDB(t), "test")

		owned, err := factory.NewBase(t.Context())
		require.NoError(t, err)

		lentCtx := owned.Lend()

		borrowed, err := factory.NewBorrowedBase(lentCtx)
		require.NoError(t, err)
		assert.Equal(t, owned.IDB(), borrowed.IDB())
	})

	t.Run("lend does not affect the original context", func(t *testing.T) {
		factory := pguowbase.NewFactory(newTestDB(t), "test")

		originalCtx := t.Context()
		owned, err := factory.NewBase(originalCtx)
		require.NoError(t, err)

		_ = owned.Lend()

		_, err = factory.NewBorrowedBase(originalCtx)
		require.Error(t, err)
	})
}

func TestApplyChanges(t *testing.T) {
	t.Run("commits an owned UOW", func(t *testing.T) {
		factory := pguowbase.NewFactory(newTestDB(t), "test")

		base, err := factory.NewBase(t.Context())
		require.NoError(t, err)

		err = base.ApplyChanges()
		assert.NoError(t, err)
	})

	t.Run("returns error on borrowed UOW", func(t *testing.T) {
		factory := pguowbase.NewFactory(newTestDB(t), "test")

		owned, err := factory.NewBase(t.Context())
		require.NoError(t, err)

		borrowed, err := factory.NewBorrowedBase(owned.Lend())
		require.NoError(t, err)

		err = borrowed.ApplyChanges()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot commit a borrowed UOW")
	})
}

func TestDiscardUnapplied(t *testing.T) {
	t.Run("rolls back an owned UOW", func(t *testing.T) {
		factory := pguowbase.NewFactory(newTestDB(t), "test")

		base, err := factory.NewBase(t.Context())
		require.NoError(t, err)

		assert.NotPanics(t, func() {
			base.DiscardUnapplied()
		})
	})

	t.Run("does not panic on borrowed UOW", func(t *testing.T) {
		factory := pguowbase.NewFactory(newTestDB(t), "test")

		owned, err := factory.NewBase(t.Context())
		require.NoError(t, err)

		borrowed, err := factory.NewBorrowedBase(owned.Lend())
		require.NoError(t, err)

		assert.NotPanics(t, func() {
			borrowed.DiscardUnapplied()
		})
	})

	t.Run("safe to call after ApplyChanges on owned UOW", func(t *testing.T) {
		factory := pguowbase.NewFactory(newTestDB(t), "test")

		base, err := factory.NewBase(t.Context())
		require.NoError(t, err)

		err = base.ApplyChanges()
		require.NoError(t, err)

		assert.NotPanics(t, func() {
			base.DiscardUnapplied()
		})
	})
}

func TestIDB(t *testing.T) {
	t.Run("returns non-nil IDB for owned base", func(t *testing.T) {
		factory := pguowbase.NewFactory(newTestDB(t), "test")

		base, err := factory.NewBase(t.Context())
		require.NoError(t, err)

		assert.NotNil(t, base.IDB())
	})

	t.Run("returns non-nil IDB for borrowed base", func(t *testing.T) {
		factory := pguowbase.NewFactory(newTestDB(t), "test")

		owned, err := factory.NewBase(t.Context())
		require.NoError(t, err)

		borrowed, err := factory.NewBorrowedBase(owned.Lend())
		require.NoError(t, err)

		assert.NotNil(t, borrowed.IDB())
	})
}

// testUOW is a minimal UOW implementation for testing GenericFactory.
type testUOW struct {
	*pguowbase.Base
}

func TestGenericFactory(t *testing.T) {
	constructor := func(base *pguowbase.Base) *testUOW {
		return &testUOW{Base: base}
	}

	t.Run("NewUOW creates owned UOW via constructor", func(t *testing.T) {
		factory := pguowbase.NewGenericFactory(newTestDB(t), "test", constructor)

		uow, err := factory.NewUOW(t.Context())
		require.NoError(t, err)
		assert.NotNil(t, uow)
		assert.NotNil(t, uow.IDB())
	})

	t.Run("NewBorrowed creates borrowed UOW via constructor", func(t *testing.T) {
		factory := pguowbase.NewGenericFactory(newTestDB(t), "test", constructor)

		owned, err := factory.NewUOW(t.Context())
		require.NoError(t, err)

		borrowed, err := factory.NewBorrowed(owned.Lend())
		require.NoError(t, err)
		assert.NotNil(t, borrowed)
		assert.Equal(t, owned.IDB(), borrowed.IDB())
	})

	t.Run("NewBorrowed returns error when no transaction lent", func(t *testing.T) {
		factory := pguowbase.NewGenericFactory(newTestDB(t), "test", constructor)

		uow, err := factory.NewBorrowed(t.Context())
		require.Error(t, err)
		assert.Nil(t, uow)
	})

	t.Run("satisfies uowbase.Factory interface", func(t *testing.T) {
		factory := pguowbase.NewGenericFactory(newTestDB(t), "test", constructor)

		var _ uowbase.Factory[*testUOW] = factory
	})
}
