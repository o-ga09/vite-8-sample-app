package repository

import (
	"context"

	"github.com/google/uuid"
	infradb "github.com/o-ga09/vite-8-sample-app/internal/infra/db"
	"github.com/o-ga09/vite-8-sample-app/internal/infra/dbgen"
	"github.com/stephenafamo/bob"
)

type WorkspaceRepository interface {
	Create(ctx context.Context, exec bob.Executor, setter *dbgen.WorkspaceSetter) (*dbgen.Workspace, error)
	Get(ctx context.Context, exec bob.Executor, id uuid.UUID) (*dbgen.Workspace, error)
	List(ctx context.Context, exec bob.Executor) (dbgen.WorkspaceSlice, error)
	Update(ctx context.Context, exec bob.Executor, id uuid.UUID, setter *dbgen.WorkspaceSetter) (*dbgen.Workspace, error)
	Delete(ctx context.Context, exec bob.Executor, id uuid.UUID) error
}

type workspaceRepository struct{}

func NewWorkspaceRepository() WorkspaceRepository {
	return &workspaceRepository{}
}

func (r *workspaceRepository) Create(ctx context.Context, exec bob.Executor, setter *dbgen.WorkspaceSetter) (*dbgen.Workspace, error) {
	ctx = infradb.GlobalExec(ctx)
	return dbgen.Workspaces.Insert(setter).One(ctx, exec)
}

func (r *workspaceRepository) Get(ctx context.Context, exec bob.Executor, id uuid.UUID) (*dbgen.Workspace, error) {
	ctx = infradb.GlobalExec(ctx)
	return dbgen.Workspaces.Query(
		dbgen.SelectWhere.Workspaces.ID.EQ(id),
	).One(ctx, exec)
}

func (r *workspaceRepository) List(ctx context.Context, exec bob.Executor) (dbgen.WorkspaceSlice, error) {
	ctx = infradb.GlobalExec(ctx)
	return dbgen.Workspaces.Query().All(ctx, exec)
}

func (r *workspaceRepository) Update(ctx context.Context, exec bob.Executor, id uuid.UUID, setter *dbgen.WorkspaceSetter) (*dbgen.Workspace, error) {
	ctx = infradb.GlobalExec(ctx)
	return dbgen.Workspaces.Update(
		setter.UpdateMod(),
		dbgen.UpdateWhere.Workspaces.ID.EQ(id),
	).One(ctx, exec)
}

func (r *workspaceRepository) Delete(ctx context.Context, exec bob.Executor, id uuid.UUID) error {
	ctx = infradb.GlobalExec(ctx)
	_, err := dbgen.Workspaces.Delete(
		dbgen.DeleteWhere.Workspaces.ID.EQ(id),
	).Exec(ctx, exec)
	return err
}
