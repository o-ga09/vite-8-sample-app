package db_test

import (
	"context"
	"errors"
	"testing"

	infradb "github.com/o-ga09/vite-8-sample-app/internal/infra/db"
)

func TestWithTx_CommitsOnSuccess(t *testing.T) {
	committed := false
	rolled := false

	tx := &fakeTx{
		commitFn:   func() error { committed = true; return nil },
		rollbackFn: func() error { rolled = true; return nil },
	}

	err := infradb.WithTx(context.Background(), tx, func(_ context.Context) error {
		return nil
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !committed {
		t.Error("expected Commit to be called")
	}
	if rolled {
		t.Error("expected Rollback NOT to be called on success")
	}
}

func TestWithTx_RollsBackOnError(t *testing.T) {
	committed := false
	rolled := false

	tx := &fakeTx{
		commitFn:   func() error { committed = true; return nil },
		rollbackFn: func() error { rolled = true; return nil },
	}

	sentinel := errors.New("fn error")
	err := infradb.WithTx(context.Background(), tx, func(_ context.Context) error {
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
	if committed {
		t.Error("expected Commit NOT to be called on error")
	}
	if !rolled {
		t.Error("expected Rollback to be called on error")
	}
}

func TestWithTx_RollsBackOnPanic(t *testing.T) {
	rolled := false

	tx := &fakeTx{
		commitFn:   func() error { return nil },
		rollbackFn: func() error { rolled = true; return nil },
	}

	defer func() {
		r := recover()
		if r == nil {
			t.Error("expected panic to propagate")
		}
		if !rolled {
			t.Error("expected Rollback to be called on panic")
		}
	}()

	_ = infradb.WithTx(context.Background(), tx, func(_ context.Context) error {
		panic("test panic")
	})
}

// fakeTx implements infra/db.Tx for testing.
type fakeTx struct {
	commitFn   func() error
	rollbackFn func() error
}

func (f *fakeTx) Commit() error   { return f.commitFn() }
func (f *fakeTx) Rollback() error { return f.rollbackFn() }
