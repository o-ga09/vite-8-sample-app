package repository

import (
	"context"

	"github.com/google/uuid"
	infradb "github.com/o-ga09/vite-8-sample-app/internal/infra/db"
	"github.com/o-ga09/vite-8-sample-app/internal/infra/dbgen"
	"github.com/stephenafamo/bob"
)

type CategoryRepository interface {
	Create(ctx context.Context, exec bob.Executor, setter *dbgen.CategorySetter) (*dbgen.Category, error)
	Get(ctx context.Context, exec bob.Executor, workspaceID, id uuid.UUID) (*dbgen.Category, error)
	List(ctx context.Context, exec bob.Executor, workspaceID uuid.UUID) (dbgen.CategorySlice, error)
	Update(ctx context.Context, exec bob.Executor, workspaceID, id uuid.UUID, setter *dbgen.CategorySetter) (*dbgen.Category, error)
	Delete(ctx context.Context, exec bob.Executor, workspaceID, id uuid.UUID) error
}

type categoryRepository struct{}

func NewCategoryRepository() CategoryRepository {
	return &categoryRepository{}
}

func (r *categoryRepository) Create(ctx context.Context, exec bob.Executor, setter *dbgen.CategorySetter) (*dbgen.Category, error) {
	wsID, ok := infradb.WorkspaceIDFromContext(ctx)
	if ok {
		ctx = infradb.WorkspaceScopedExec(ctx, wsID)
	}
	return dbgen.Categories.Insert(setter).One(ctx, exec)
}

func (r *categoryRepository) Get(ctx context.Context, exec bob.Executor, workspaceID, id uuid.UUID) (*dbgen.Category, error) {
	ctx = infradb.WorkspaceScopedExec(ctx, workspaceID)
	return dbgen.Categories.Query(
		dbgen.SelectWhere.Categories.WorkspaceID.EQ(workspaceID),
		dbgen.SelectWhere.Categories.ID.EQ(id),
	).One(ctx, exec)
}

func (r *categoryRepository) List(ctx context.Context, exec bob.Executor, workspaceID uuid.UUID) (dbgen.CategorySlice, error) {
	ctx = infradb.WorkspaceScopedExec(ctx, workspaceID)
	return dbgen.Categories.Query(
		dbgen.SelectWhere.Categories.WorkspaceID.EQ(workspaceID),
	).All(ctx, exec)
}

func (r *categoryRepository) Update(ctx context.Context, exec bob.Executor, workspaceID, id uuid.UUID, setter *dbgen.CategorySetter) (*dbgen.Category, error) {
	ctx = infradb.WorkspaceScopedExec(ctx, workspaceID)
	return dbgen.Categories.Update(
		setter.UpdateMod(),
		dbgen.UpdateWhere.Categories.WorkspaceID.EQ(workspaceID),
		dbgen.UpdateWhere.Categories.ID.EQ(id),
	).One(ctx, exec)
}

func (r *categoryRepository) Delete(ctx context.Context, exec bob.Executor, workspaceID, id uuid.UUID) error {
	ctx = infradb.WorkspaceScopedExec(ctx, workspaceID)
	_, err := dbgen.Categories.Delete(
		dbgen.DeleteWhere.Categories.WorkspaceID.EQ(workspaceID),
		dbgen.DeleteWhere.Categories.ID.EQ(id),
	).Exec(ctx, exec)
	return err
}
