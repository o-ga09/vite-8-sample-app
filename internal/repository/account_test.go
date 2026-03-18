package repository_test

import (
	"context"
	"testing"

	"github.com/aarondl/opt/omit"
	"github.com/google/uuid"
	"github.com/o-ga09/vite-8-sample-app/internal/infra/dbgen"
	dbenums "github.com/o-ga09/vite-8-sample-app/internal/infra/dbgen/dbenums"
	"github.com/o-ga09/vite-8-sample-app/internal/infra/dbgen/factory"
	"github.com/o-ga09/vite-8-sample-app/internal/repository"
	"github.com/shopspring/decimal"
)

func TestNewAccountRepository_ReturnsNonNil(t *testing.T) {
	r := repository.NewAccountRepository()
	if r == nil {
		t.Fatal("expected non-nil AccountRepository")
	}
}

func TestAccountRepository_Create(t *testing.T) {
	if testDB == nil {
		t.Skip("skipping: no DB connection")
	}

	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)

	tx, err := testDB.Begin(ctx)
	if err != nil {
		t.Fatalf("failed to begin tx: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	ws, err := factory.New().NewWorkspace().Create(ctx, tx)
	if err != nil {
		t.Fatalf("factory workspace create failed: %v", err)
	}

	r := repository.NewAccountRepository()
	setter := &dbgen.AccountSetter{
		ID:             omit.From(uuid.New()),
		WorkspaceID:    omit.From(ws.ID),
		Name:           omit.From("test-account"),
		AccountType:    omit.From(dbenums.AccountTypeCash),
		InitialBalance: omit.From(decimal.NewFromInt(0)),
	}

	acc, err := r.Create(ctx, tx, setter)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if acc == nil {
		t.Fatal("expected non-nil account")
	}
	if acc.WorkspaceID != ws.ID {
		t.Errorf("got WorkspaceID %v, want %v", acc.WorkspaceID, ws.ID)
	}
}

func TestAccountRepository_Get_ExistingAccount(t *testing.T) {
	if testDB == nil {
		t.Skip("skipping: no DB connection")
	}

	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)

	tx, err := testDB.Begin(ctx)
	if err != nil {
		t.Fatalf("failed to begin tx: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	acc, err := factory.New().NewAccount().Create(ctx, tx)
	if err != nil {
		t.Fatalf("factory.Create failed: %v", err)
	}

	r := repository.NewAccountRepository()
	got, err := r.Get(ctx, tx, acc.WorkspaceID, acc.ID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got.ID != acc.ID {
		t.Errorf("got ID %v, want %v", got.ID, acc.ID)
	}
}

func TestAccountRepository_Get_NonExistentID(t *testing.T) {
	if testDB == nil {
		t.Skip("skipping: no DB connection")
	}

	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)

	tx, err := testDB.Begin(ctx)
	if err != nil {
		t.Fatalf("failed to begin tx: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	r := repository.NewAccountRepository()
	_, err = r.Get(ctx, tx, uuid.New(), uuid.New())
	if err == nil {
		t.Fatal("expected error for non-existent account, got nil")
	}
}

func TestAccountRepository_List(t *testing.T) {
	if testDB == nil {
		t.Skip("skipping: no DB connection")
	}

	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)

	tx, err := testDB.Begin(ctx)
	if err != nil {
		t.Fatalf("failed to begin tx: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Create workspace once, then create 2 accounts within it.
	ws, err := factory.New().NewWorkspace().Create(ctx, tx)
	if err != nil {
		t.Fatalf("factory workspace create failed: %v", err)
	}

	f := factory.New()
	for i := 0; i < 2; i++ {
		if _, createErr := f.NewAccount(factory.AccountMods.WithExistingWorkspace(ws)).Create(ctx, tx); createErr != nil {
			t.Fatalf("factory account create(%d) failed: %v", i, createErr)
		}
	}

	r := repository.NewAccountRepository()
	list, err := r.List(ctx, tx, ws.ID)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(list) < 2 {
		t.Errorf("expected at least 2 accounts, got %d", len(list))
	}
}

func TestAccountRepository_Update(t *testing.T) {
	if testDB == nil {
		t.Skip("skipping: no DB connection")
	}

	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)

	tx, err := testDB.Begin(ctx)
	if err != nil {
		t.Fatalf("failed to begin tx: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	acc, err := factory.New().NewAccount().Create(ctx, tx)
	if err != nil {
		t.Fatalf("factory.Create failed: %v", err)
	}

	r := repository.NewAccountRepository()
	const newName = "updated-account"
	setter := &dbgen.AccountSetter{
		Name: omit.From(newName),
	}

	updated, err := r.Update(ctx, tx, acc.WorkspaceID, acc.ID, setter)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if updated.Name != newName {
		t.Errorf("got Name %q, want %q", updated.Name, newName)
	}
}

func TestAccountRepository_Delete(t *testing.T) {
	if testDB == nil {
		t.Skip("skipping: no DB connection")
	}

	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)

	tx, err := testDB.Begin(ctx)
	if err != nil {
		t.Fatalf("failed to begin tx: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	acc, err := factory.New().NewAccount().Create(ctx, tx)
	if err != nil {
		t.Fatalf("factory.Create failed: %v", err)
	}

	r := repository.NewAccountRepository()
	if err = r.Delete(ctx, tx, acc.WorkspaceID, acc.ID); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = r.Get(ctx, tx, acc.WorkspaceID, acc.ID)
	if err == nil {
		t.Fatal("expected error after delete, got nil")
	}
}
