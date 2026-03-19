package service

import (
	"context"

	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/google/uuid"
	"github.com/o-ga09/vite-8-sample-app/internal/infra/dbgen"
	enums "github.com/o-ga09/vite-8-sample-app/internal/infra/dbgen/dbenums"
	"github.com/o-ga09/vite-8-sample-app/internal/repository"
	"github.com/stephenafamo/bob"
)

// CategoryService handles category business logic.
type CategoryService interface {
	Create(ctx context.Context, workspaceID uuid.UUID, name string, categoryType enums.CategoryType, parentID *uuid.UUID) (*dbgen.Category, error)
	Get(ctx context.Context, workspaceID, id uuid.UUID) (*dbgen.Category, error)
	List(ctx context.Context, workspaceID uuid.UUID) (dbgen.CategorySlice, error)
	Update(ctx context.Context, workspaceID, id uuid.UUID, name string, categoryType enums.CategoryType, parentID *uuid.UUID) (*dbgen.Category, error)
	Delete(ctx context.Context, workspaceID, id uuid.UUID) error
}

type categoryService struct {
	db   bob.DB
	repo repository.CategoryRepository
}

// NewCategoryService creates a new CategoryService.
func NewCategoryService(db bob.DB, repo repository.CategoryRepository) CategoryService {
	return &categoryService{db: db, repo: repo}
}

func (s *categoryService) Create(ctx context.Context, workspaceID uuid.UUID, name string, categoryType enums.CategoryType, parentID *uuid.UUID) (*dbgen.Category, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	setter := &dbgen.CategorySetter{
		ID:           omit.From(id),
		WorkspaceID:  omit.From(workspaceID),
		Name:         omit.From(name),
		CategoryType: omit.From(categoryType),
	}
	if parentID != nil {
		setter.ParentID = omitnull.From(*parentID)
	}
	return s.repo.Create(ctx, s.db, setter)
}

func (s *categoryService) Get(ctx context.Context, workspaceID, id uuid.UUID) (*dbgen.Category, error) {
	c, err := s.repo.Get(ctx, s.db, workspaceID, id)
	if err != nil {
		return nil, mapRepoErr(err)
	}
	return c, nil
}

func (s *categoryService) List(ctx context.Context, workspaceID uuid.UUID) (dbgen.CategorySlice, error) {
	return s.repo.List(ctx, s.db, workspaceID)
}

func (s *categoryService) Update(ctx context.Context, workspaceID, id uuid.UUID, name string, categoryType enums.CategoryType, parentID *uuid.UUID) (*dbgen.Category, error) {
	setter := &dbgen.CategorySetter{
		Name:         omit.From(name),
		CategoryType: omit.From(categoryType),
	}
	if parentID != nil {
		setter.ParentID = omitnull.From(*parentID)
	}
	c, err := s.repo.Update(ctx, s.db, workspaceID, id, setter)
	if err != nil {
		return nil, mapRepoErr(err)
	}
	return c, nil
}

func (s *categoryService) Delete(ctx context.Context, workspaceID, id uuid.UUID) error {
	err := s.repo.Delete(ctx, s.db, workspaceID, id)
	if err != nil {
		return mapRepoErr(err)
	}
	return nil
}
