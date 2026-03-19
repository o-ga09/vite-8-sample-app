package repository

import (
	"context"

	"github.com/google/uuid"
	infradb "github.com/o-ga09/vite-8-sample-app/internal/infra/db"
	"github.com/o-ga09/vite-8-sample-app/internal/infra/dbgen"
	"github.com/stephenafamo/bob"
)

type MemberRepository interface {
	Create(ctx context.Context, exec bob.Executor, setter *dbgen.MemberSetter) (*dbgen.Member, error)
	Get(ctx context.Context, exec bob.Executor, workspaceID, id uuid.UUID) (*dbgen.Member, error)
	List(ctx context.Context, exec bob.Executor, workspaceID uuid.UUID) (dbgen.MemberSlice, error)
	Update(ctx context.Context, exec bob.Executor, workspaceID, id uuid.UUID, setter *dbgen.MemberSetter) (*dbgen.Member, error)
	Delete(ctx context.Context, exec bob.Executor, workspaceID, id uuid.UUID) error
}

type memberRepository struct{}

func NewMemberRepository() MemberRepository {
	return &memberRepository{}
}

func (r *memberRepository) Create(ctx context.Context, exec bob.Executor, setter *dbgen.MemberSetter) (*dbgen.Member, error) {
	wsID, ok := infradb.WorkspaceIDFromContext(ctx)
	if ok {
		ctx = infradb.WorkspaceScopedExec(ctx, wsID)
	}
	return dbgen.Members.Insert(setter).One(ctx, exec)
}

func (r *memberRepository) Get(ctx context.Context, exec bob.Executor, workspaceID, id uuid.UUID) (*dbgen.Member, error) {
	ctx = infradb.WorkspaceScopedExec(ctx, workspaceID)
	return dbgen.Members.Query(
		dbgen.SelectWhere.Members.WorkspaceID.EQ(workspaceID),
		dbgen.SelectWhere.Members.ID.EQ(id),
	).One(ctx, exec)
}

func (r *memberRepository) List(ctx context.Context, exec bob.Executor, workspaceID uuid.UUID) (dbgen.MemberSlice, error) {
	ctx = infradb.WorkspaceScopedExec(ctx, workspaceID)
	return dbgen.Members.Query(
		dbgen.SelectWhere.Members.WorkspaceID.EQ(workspaceID),
	).All(ctx, exec)
}

func (r *memberRepository) Update(ctx context.Context, exec bob.Executor, workspaceID, id uuid.UUID, setter *dbgen.MemberSetter) (*dbgen.Member, error) {
	ctx = infradb.WorkspaceScopedExec(ctx, workspaceID)
	return dbgen.Members.Update(
		setter.UpdateMod(),
		dbgen.UpdateWhere.Members.WorkspaceID.EQ(workspaceID),
		dbgen.UpdateWhere.Members.ID.EQ(id),
	).One(ctx, exec)
}

func (r *memberRepository) Delete(ctx context.Context, exec bob.Executor, workspaceID, id uuid.UUID) error {
	ctx = infradb.WorkspaceScopedExec(ctx, workspaceID)
	_, err := dbgen.Members.Delete(
		dbgen.DeleteWhere.Members.WorkspaceID.EQ(workspaceID),
		dbgen.DeleteWhere.Members.ID.EQ(id),
	).Exec(ctx, exec)
	return err
}
