package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	infradb "github.com/o-ga09/vite-8-sample-app/internal/infra/db"
	"github.com/o-ga09/vite-8-sample-app/internal/infra/dbgen"
	"github.com/stephenafamo/bob"
)

type TransactionRepository interface {
	Create(ctx context.Context, exec bob.Executor, setter *dbgen.TransactionSetter) (*dbgen.Transaction, error)
	Get(ctx context.Context, exec bob.Executor, workspaceID, id uuid.UUID) (*dbgen.Transaction, error)
	List(ctx context.Context, exec bob.Executor, workspaceID uuid.UUID) (dbgen.TransactionSlice, error)
	ListByPeriod(ctx context.Context, exec bob.Executor, workspaceID uuid.UUID, from, to time.Time) (dbgen.TransactionSlice, error)
	Update(ctx context.Context, exec bob.Executor, workspaceID, id uuid.UUID, setter *dbgen.TransactionSetter) (*dbgen.Transaction, error)
	Delete(ctx context.Context, exec bob.Executor, workspaceID, id uuid.UUID) error
}

type transactionRepository struct{}

func NewTransactionRepository() TransactionRepository {
	return &transactionRepository{}
}

func (r *transactionRepository) Create(ctx context.Context, exec bob.Executor, setter *dbgen.TransactionSetter) (*dbgen.Transaction, error) {
	wsID, ok := infradb.WorkspaceIDFromContext(ctx)
	if ok {
		ctx = infradb.WorkspaceScopedExec(ctx, wsID)
	}
	return dbgen.Transactions.Insert(setter).One(ctx, exec)
}

func (r *transactionRepository) Get(ctx context.Context, exec bob.Executor, workspaceID, id uuid.UUID) (*dbgen.Transaction, error) {
	ctx = infradb.WorkspaceScopedExec(ctx, workspaceID)
	return dbgen.Transactions.Query(
		dbgen.SelectWhere.Transactions.WorkspaceID.EQ(workspaceID),
		dbgen.SelectWhere.Transactions.ID.EQ(id),
		dbgen.Preload.Transaction.Account(),
		dbgen.Preload.Transaction.CounterpartyAccount(),
		dbgen.Preload.Transaction.Category(),
	).One(ctx, exec)
}

func (r *transactionRepository) List(ctx context.Context, exec bob.Executor, workspaceID uuid.UUID) (dbgen.TransactionSlice, error) {
	ctx = infradb.WorkspaceScopedExec(ctx, workspaceID)
	return dbgen.Transactions.Query(
		dbgen.SelectWhere.Transactions.WorkspaceID.EQ(workspaceID),
		dbgen.Preload.Transaction.Account(),
		dbgen.Preload.Transaction.CounterpartyAccount(),
		dbgen.Preload.Transaction.Category(),
	).All(ctx, exec)
}

func (r *transactionRepository) ListByPeriod(ctx context.Context, exec bob.Executor, workspaceID uuid.UUID, from, to time.Time) (dbgen.TransactionSlice, error) {
	ctx = infradb.WorkspaceScopedExec(ctx, workspaceID)
	return dbgen.Transactions.Query(
		dbgen.SelectWhere.Transactions.WorkspaceID.EQ(workspaceID),
		dbgen.SelectWhere.Transactions.OccurredAt.GTE(from),
		dbgen.SelectWhere.Transactions.OccurredAt.LTE(to),
		dbgen.Preload.Transaction.Account(),
		dbgen.Preload.Transaction.CounterpartyAccount(),
		dbgen.Preload.Transaction.Category(),
	).All(ctx, exec)
}

func (r *transactionRepository) Update(ctx context.Context, exec bob.Executor, workspaceID, id uuid.UUID, setter *dbgen.TransactionSetter) (*dbgen.Transaction, error) {
	ctx = infradb.WorkspaceScopedExec(ctx, workspaceID)
	return dbgen.Transactions.Update(
		setter.UpdateMod(),
		dbgen.UpdateWhere.Transactions.WorkspaceID.EQ(workspaceID),
		dbgen.UpdateWhere.Transactions.ID.EQ(id),
	).One(ctx, exec)
}

func (r *transactionRepository) Delete(ctx context.Context, exec bob.Executor, workspaceID, id uuid.UUID) error {
	ctx = infradb.WorkspaceScopedExec(ctx, workspaceID)
	_, err := dbgen.Transactions.Delete(
		dbgen.DeleteWhere.Transactions.WorkspaceID.EQ(workspaceID),
		dbgen.DeleteWhere.Transactions.ID.EQ(id),
	).Exec(ctx, exec)
	return err
}
