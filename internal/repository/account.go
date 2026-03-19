package repository

import (
	"context"

	"github.com/google/uuid"
	infradb "github.com/o-ga09/vite-8-sample-app/internal/infra/db"
	"github.com/o-ga09/vite-8-sample-app/internal/infra/dbgen"
	"github.com/stephenafamo/bob"
)

type AccountRepository interface {
	Create(ctx context.Context, exec bob.Executor, setter *dbgen.AccountSetter) (*dbgen.Account, error)
	Get(ctx context.Context, exec bob.Executor, workspaceID, id uuid.UUID) (*dbgen.Account, error)
	List(ctx context.Context, exec bob.Executor, workspaceID uuid.UUID) (dbgen.AccountSlice, error)
	Update(ctx context.Context, exec bob.Executor, workspaceID, id uuid.UUID, setter *dbgen.AccountSetter) (*dbgen.Account, error)
	Delete(ctx context.Context, exec bob.Executor, workspaceID, id uuid.UUID) error
}

type accountRepository struct{}

func NewAccountRepository() AccountRepository {
	return &accountRepository{}
}

func (r *accountRepository) Create(ctx context.Context, exec bob.Executor, setter *dbgen.AccountSetter) (*dbgen.Account, error) {
	wsID, ok := infradb.WorkspaceIDFromContext(ctx)
	if ok {
		ctx = infradb.WorkspaceScopedExec(ctx, wsID)
	}
	return dbgen.Accounts.Insert(setter).One(ctx, exec)
}

func (r *accountRepository) Get(ctx context.Context, exec bob.Executor, workspaceID, id uuid.UUID) (*dbgen.Account, error) {
	ctx = infradb.WorkspaceScopedExec(ctx, workspaceID)
	return dbgen.Accounts.Query(
		dbgen.SelectWhere.Accounts.WorkspaceID.EQ(workspaceID),
		dbgen.SelectWhere.Accounts.ID.EQ(id),
	).One(ctx, exec)
}

func (r *accountRepository) List(ctx context.Context, exec bob.Executor, workspaceID uuid.UUID) (dbgen.AccountSlice, error) {
	ctx = infradb.WorkspaceScopedExec(ctx, workspaceID)
	return dbgen.Accounts.Query(
		dbgen.SelectWhere.Accounts.WorkspaceID.EQ(workspaceID),
	).All(ctx, exec)
}

func (r *accountRepository) Update(ctx context.Context, exec bob.Executor, workspaceID, id uuid.UUID, setter *dbgen.AccountSetter) (*dbgen.Account, error) {
	ctx = infradb.WorkspaceScopedExec(ctx, workspaceID)
	return dbgen.Accounts.Update(
		setter.UpdateMod(),
		dbgen.UpdateWhere.Accounts.WorkspaceID.EQ(workspaceID),
		dbgen.UpdateWhere.Accounts.ID.EQ(id),
	).One(ctx, exec)
}

func (r *accountRepository) Delete(ctx context.Context, exec bob.Executor, workspaceID, id uuid.UUID) error {
	ctx = infradb.WorkspaceScopedExec(ctx, workspaceID)
	_, err := dbgen.Accounts.Delete(
		dbgen.DeleteWhere.Accounts.WorkspaceID.EQ(workspaceID),
		dbgen.DeleteWhere.Accounts.ID.EQ(id),
	).Exec(ctx, exec)
	return err
}
