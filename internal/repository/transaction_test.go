package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/google/uuid"
	"github.com/o-ga09/vite-8-sample-app/internal/infra/dbgen"
	dbenums "github.com/o-ga09/vite-8-sample-app/internal/infra/dbgen/dbenums"
	"github.com/o-ga09/vite-8-sample-app/internal/infra/dbgen/factory"
	"github.com/o-ga09/vite-8-sample-app/internal/repository"
	"github.com/shopspring/decimal"
	"github.com/stephenafamo/bob"
)

func TestNewTransactionRepository_ReturnsNonNil(t *testing.T) {
	r := repository.NewTransactionRepository()
	if r == nil {
		t.Fatal("expected non-nil TransactionRepository")
	}
}

// createIncomeTransaction is a helper that inserts an income-type transaction
// satisfying the chk_transactions_type_relations DB constraint.
func createIncomeTransaction(ctx context.Context, t *testing.T, exec bob.Executor, wsID, accountID, catID uuid.UUID) *dbgen.Transaction {
	t.Helper()
	setter := &dbgen.TransactionSetter{
		ID:              omit.From(uuid.New()),
		WorkspaceID:     omit.From(wsID),
		TransactionType: omit.From(dbenums.TransactionTypeIncome),
		AccountID:       omitnull.From(accountID),
		CategoryID:      omitnull.From(catID),
		Amount:          omit.From(decimal.NewFromInt(1000)),
		OccurredAt:      omit.From(time.Now()),
	}
	txn, err := dbgen.Transactions.Insert(setter).One(ctx, exec)
	if err != nil {
		t.Fatalf("createIncomeTransaction: %v", err)
	}
	return txn
}

// setupWorkspaceWithAccountAndCategory creates a workspace, an account, and an
// income-category all sharing the same workspace_id. It returns their IDs.
func setupWorkspaceWithAccountAndCategory(ctx context.Context, t *testing.T, exec bob.Executor) (wsID, accountID, catID uuid.UUID) {
	t.Helper()
	f := factory.New()

	ws, err := f.NewWorkspace().Create(ctx, exec)
	if err != nil {
		t.Fatalf("create workspace: %v", err)
	}

	acc, err := f.NewAccount(factory.AccountMods.WithExistingWorkspace(ws)).Create(ctx, exec)
	if err != nil {
		t.Fatalf("create account: %v", err)
	}

	catSetter := &dbgen.CategorySetter{
		ID:           omit.From(uuid.New()),
		WorkspaceID:  omit.From(ws.ID),
		Name:         omit.From("テスト収入カテゴリ"),
		CategoryType: omit.From(dbenums.CategoryTypeIncome),
	}
	cat, err := dbgen.Categories.Insert(catSetter).One(ctx, exec)
	if err != nil {
		t.Fatalf("create category: %v", err)
	}

	return ws.ID, acc.ID, cat.ID
}

func TestTransactionRepository_Create(t *testing.T) {
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

	wsID, accountID, catID := setupWorkspaceWithAccountAndCategory(ctx, t, tx)

	r := repository.NewTransactionRepository()
	setter := &dbgen.TransactionSetter{
		ID:              omit.From(uuid.New()),
		WorkspaceID:     omit.From(wsID),
		TransactionType: omit.From(dbenums.TransactionTypeIncome),
		AccountID:       omitnull.From(accountID),
		CategoryID:      omitnull.From(catID),
		Amount:          omit.From(decimal.NewFromInt(500)),
		OccurredAt:      omit.From(time.Now()),
	}

	txn, err := r.Create(ctx, tx, setter)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if txn == nil {
		t.Fatal("expected non-nil transaction")
	}
	if txn.WorkspaceID != wsID {
		t.Errorf("got WorkspaceID %v, want %v", txn.WorkspaceID, wsID)
	}
}

func TestTransactionRepository_Get_ExistingTransaction(t *testing.T) {
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

	wsID, accountID, catID := setupWorkspaceWithAccountAndCategory(ctx, t, tx)
	created := createIncomeTransaction(ctx, t, tx, wsID, accountID, catID)

	r := repository.NewTransactionRepository()
	got, err := r.Get(ctx, tx, wsID, created.ID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got.ID != created.ID {
		t.Errorf("got ID %v, want %v", got.ID, created.ID)
	}
}

func TestTransactionRepository_Get_NonExistentID(t *testing.T) {
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

	r := repository.NewTransactionRepository()
	_, err = r.Get(ctx, tx, uuid.New(), uuid.New())
	if err == nil {
		t.Fatal("expected error for non-existent transaction, got nil")
	}
}

func TestTransactionRepository_List(t *testing.T) {
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

	wsID, accountID, catID := setupWorkspaceWithAccountAndCategory(ctx, t, tx)
	createIncomeTransaction(ctx, t, tx, wsID, accountID, catID)
	createIncomeTransaction(ctx, t, tx, wsID, accountID, catID)

	r := repository.NewTransactionRepository()
	list, err := r.List(ctx, tx, wsID)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(list) < 2 {
		t.Errorf("expected at least 2 transactions, got %d", len(list))
	}
}

func TestTransactionRepository_ListByPeriod(t *testing.T) {
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

	wsID, accountID, catID := setupWorkspaceWithAccountAndCategory(ctx, t, tx)

	now := time.Now()
	// Create a transaction within the period and one outside.
	setterIn := &dbgen.TransactionSetter{
		ID:              omit.From(uuid.New()),
		WorkspaceID:     omit.From(wsID),
		TransactionType: omit.From(dbenums.TransactionTypeIncome),
		AccountID:       omitnull.From(accountID),
		CategoryID:      omitnull.From(catID),
		Amount:          omit.From(decimal.NewFromInt(100)),
		OccurredAt:      omit.From(now),
	}
	setterOut := &dbgen.TransactionSetter{
		ID:              omit.From(uuid.New()),
		WorkspaceID:     omit.From(wsID),
		TransactionType: omit.From(dbenums.TransactionTypeIncome),
		AccountID:       omitnull.From(accountID),
		CategoryID:      omitnull.From(catID),
		Amount:          omit.From(decimal.NewFromInt(200)),
		OccurredAt:      omit.From(now.Add(-30 * 24 * time.Hour)),
	}

	if _, err := dbgen.Transactions.Insert(setterIn).One(ctx, tx); err != nil {
		t.Fatalf("insert in-period tx: %v", err)
	}
	if _, err := dbgen.Transactions.Insert(setterOut).One(ctx, tx); err != nil {
		t.Fatalf("insert out-period tx: %v", err)
	}

	r := repository.NewTransactionRepository()
	from := now.Add(-24 * time.Hour)
	to := now.Add(24 * time.Hour)
	list, err := r.ListByPeriod(ctx, tx, wsID, from, to)
	if err != nil {
		t.Fatalf("ListByPeriod failed: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 transaction in period, got %d", len(list))
	}
	if list[0].OccurredAt.Before(from) || list[0].OccurredAt.After(to) {
		t.Errorf("transaction is outside the requested period")
	}
}

func TestTransactionRepository_Update(t *testing.T) {
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

	wsID, accountID, catID := setupWorkspaceWithAccountAndCategory(ctx, t, tx)
	created := createIncomeTransaction(ctx, t, tx, wsID, accountID, catID)

	r := repository.NewTransactionRepository()
	newAmount := decimal.NewFromInt(9999)
	setter := &dbgen.TransactionSetter{
		Amount: omit.From(newAmount),
	}

	updated, err := r.Update(ctx, tx, wsID, created.ID, setter)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if !updated.Amount.Equal(newAmount) {
		t.Errorf("got Amount %v, want %v", updated.Amount, newAmount)
	}
}

func TestTransactionRepository_Delete(t *testing.T) {
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

	wsID, accountID, catID := setupWorkspaceWithAccountAndCategory(ctx, t, tx)
	created := createIncomeTransaction(ctx, t, tx, wsID, accountID, catID)

	r := repository.NewTransactionRepository()
	if err := r.Delete(ctx, tx, wsID, created.ID); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = r.Get(ctx, tx, wsID, created.ID)
	if err == nil {
		t.Fatal("expected error after delete, got nil")
	}
}
