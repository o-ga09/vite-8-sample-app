package service

import (
	"context"

	"github.com/aarondl/opt/omit"
	"github.com/google/uuid"
	"github.com/o-ga09/vite-8-sample-app/internal/infra/dbgen"
	"github.com/o-ga09/vite-8-sample-app/internal/repository"
	"github.com/stephenafamo/bob"
)

// WorkspaceService handles workspace business logic.
type WorkspaceService interface {
	Create(ctx context.Context, name string) (*dbgen.Workspace, error)
	Get(ctx context.Context, id uuid.UUID) (*dbgen.Workspace, error)
	List(ctx context.Context) (dbgen.WorkspaceSlice, error)
	Update(ctx context.Context, id uuid.UUID, name string) (*dbgen.Workspace, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type workspaceService struct {
	db   bob.DB
	repo repository.WorkspaceRepository
}

// NewWorkspaceService creates a new WorkspaceService.
func NewWorkspaceService(db bob.DB, repo repository.WorkspaceRepository) WorkspaceService {
	return &workspaceService{db: db, repo: repo}
}

func (s *workspaceService) Create(ctx context.Context, name string) (*dbgen.Workspace, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	setter := &dbgen.WorkspaceSetter{
		ID:   omit.From(id),
		Name: omit.From(name),
	}
	return s.repo.Create(ctx, s.db, setter)
}

func (s *workspaceService) Get(ctx context.Context, id uuid.UUID) (*dbgen.Workspace, error) {
	ws, err := s.repo.Get(ctx, s.db, id)
	if err != nil {
		return nil, mapRepoErr(err)
	}
	return ws, nil
}

func (s *workspaceService) List(ctx context.Context) (dbgen.WorkspaceSlice, error) {
	return s.repo.List(ctx, s.db)
}

func (s *workspaceService) Update(ctx context.Context, id uuid.UUID, name string) (*dbgen.Workspace, error) {
	setter := &dbgen.WorkspaceSetter{
		Name: omit.From(name),
	}
	ws, err := s.repo.Update(ctx, s.db, id, setter)
	if err != nil {
		return nil, mapRepoErr(err)
	}
	return ws, nil
}

func (s *workspaceService) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.repo.Delete(ctx, s.db, id)
	if err != nil {
		return mapRepoErr(err)
	}
	return nil
}
