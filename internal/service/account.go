package service

import (
	"context"

	"github.com/aarondl/opt/omit"
	"github.com/google/uuid"
	"github.com/o-ga09/vite-8-sample-app/internal/infra/dbgen"
	enums "github.com/o-ga09/vite-8-sample-app/internal/infra/dbgen/dbenums"
	"github.com/o-ga09/vite-8-sample-app/internal/repository"
	"github.com/shopspring/decimal"
	"github.com/stephenafamo/bob"
)

// AccountService handles account business logic.
type AccountService interface {
	Create(ctx context.Context, workspaceID uuid.UUID, name string, accountType enums.AccountType, initialBalance decimal.Decimal) (*dbgen.Account, error)
	Get(ctx context.Context, workspaceID, id uuid.UUID) (*dbgen.Account, error)
	List(ctx context.Context, workspaceID uuid.UUID) (dbgen.AccountSlice, error)
	Update(ctx context.Context, workspaceID, id uuid.UUID, name string, accountType enums.AccountType) (*dbgen.Account, error)
	Delete(ctx context.Context, workspaceID, id uuid.UUID) error
}

type accountService struct {
	db   bob.DB
	repo repository.AccountRepository
}

// NewAccountService creates a new AccountService.
func NewAccountService(db bob.DB, repo repository.AccountRepository) AccountService {
	return &accountService{db: db, repo: repo}
}

func (s *accountService) Create(ctx context.Context, workspaceID uuid.UUID, name string, accountType enums.AccountType, initialBalance decimal.Decimal) (*dbgen.Account, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	setter := &dbgen.AccountSetter{
		ID:             omit.From(id),
		WorkspaceID:    omit.From(workspaceID),
		Name:           omit.From(name),
		AccountType:    omit.From(accountType),
		InitialBalance: omit.From(initialBalance),
	}
	return s.repo.Create(ctx, s.db, setter)
}

func (s *accountService) Get(ctx context.Context, workspaceID, id uuid.UUID) (*dbgen.Account, error) {
	a, err := s.repo.Get(ctx, s.db, workspaceID, id)
	if err != nil {
		return nil, mapRepoErr(err)
	}
	return a, nil
}

func (s *accountService) List(ctx context.Context, workspaceID uuid.UUID) (dbgen.AccountSlice, error) {
	return s.repo.List(ctx, s.db, workspaceID)
}

func (s *accountService) Update(ctx context.Context, workspaceID, id uuid.UUID, name string, accountType enums.AccountType) (*dbgen.Account, error) {
	setter := &dbgen.AccountSetter{
		Name:        omit.From(name),
		AccountType: omit.From(accountType),
	}
	a, err := s.repo.Update(ctx, s.db, workspaceID, id, setter)
	if err != nil {
		return nil, mapRepoErr(err)
	}
	return a, nil
}

func (s *accountService) Delete(ctx context.Context, workspaceID, id uuid.UUID) error {
	err := s.repo.Delete(ctx, s.db, workspaceID, id)
	if err != nil {
		return mapRepoErr(err)
	}
	return nil
}
