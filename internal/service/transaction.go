package service

import (
	"context"
	"time"

	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/google/uuid"
	"github.com/o-ga09/vite-8-sample-app/internal/infra/dbgen"
	enums "github.com/o-ga09/vite-8-sample-app/internal/infra/dbgen/dbenums"
	"github.com/o-ga09/vite-8-sample-app/internal/repository"
	"github.com/shopspring/decimal"
	"github.com/stephenafamo/bob"
)

// CreateTransactionInput holds the input data for creating a transaction.
type CreateTransactionInput struct {
	TransactionType       enums.TransactionType
	AccountID             *uuid.UUID
	CounterpartyAccountID *uuid.UUID
	CategoryID            *uuid.UUID
	Amount                decimal.Decimal
	OccurredAt            time.Time
	Description           *string
}

// TransactionService handles transaction business logic.
type TransactionService interface {
	Create(ctx context.Context, workspaceID uuid.UUID, input CreateTransactionInput) (*dbgen.Transaction, error)
	Get(ctx context.Context, workspaceID, id uuid.UUID) (*dbgen.Transaction, error)
	List(ctx context.Context, workspaceID uuid.UUID) (dbgen.TransactionSlice, error)
	ListByPeriod(ctx context.Context, workspaceID uuid.UUID, from, to time.Time) (dbgen.TransactionSlice, error)
	Update(ctx context.Context, workspaceID, id uuid.UUID, input CreateTransactionInput) (*dbgen.Transaction, error)
	Delete(ctx context.Context, workspaceID, id uuid.UUID) error
}

type transactionService struct {
	db   bob.DB
	repo repository.TransactionRepository
}

// NewTransactionService creates a new TransactionService.
func NewTransactionService(db bob.DB, repo repository.TransactionRepository) TransactionService {
	return &transactionService{db: db, repo: repo}
}

func (s *transactionService) Create(ctx context.Context, workspaceID uuid.UUID, input CreateTransactionInput) (*dbgen.Transaction, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	setter := buildTransactionSetter(id, workspaceID, input)
	return s.repo.Create(ctx, s.db, setter)
}

func (s *transactionService) Get(ctx context.Context, workspaceID, id uuid.UUID) (*dbgen.Transaction, error) {
	tx, err := s.repo.Get(ctx, s.db, workspaceID, id)
	if err != nil {
		return nil, mapRepoErr(err)
	}
	return tx, nil
}

func (s *transactionService) List(ctx context.Context, workspaceID uuid.UUID) (dbgen.TransactionSlice, error) {
	return s.repo.List(ctx, s.db, workspaceID)
}

func (s *transactionService) ListByPeriod(ctx context.Context, workspaceID uuid.UUID, from, to time.Time) (dbgen.TransactionSlice, error) {
	return s.repo.ListByPeriod(ctx, s.db, workspaceID, from, to)
}

func (s *transactionService) Update(ctx context.Context, workspaceID, id uuid.UUID, input CreateTransactionInput) (*dbgen.Transaction, error) {
	setter := buildTransactionSetter(uuid.UUID{}, workspaceID, input)
	tx, err := s.repo.Update(ctx, s.db, workspaceID, id, setter)
	if err != nil {
		return nil, mapRepoErr(err)
	}
	return tx, nil
}

func (s *transactionService) Delete(ctx context.Context, workspaceID, id uuid.UUID) error {
	err := s.repo.Delete(ctx, s.db, workspaceID, id)
	if err != nil {
		return mapRepoErr(err)
	}
	return nil
}

func buildTransactionSetter(id, workspaceID uuid.UUID, input CreateTransactionInput) *dbgen.TransactionSetter {
	setter := &dbgen.TransactionSetter{
		WorkspaceID:     omit.From(workspaceID),
		TransactionType: omit.From(input.TransactionType),
		Amount:          omit.From(input.Amount),
		OccurredAt:      omit.From(input.OccurredAt),
	}
	if id != (uuid.UUID{}) {
		setter.ID = omit.From(id)
	}
	if input.AccountID != nil {
		setter.AccountID = omitnull.From(*input.AccountID)
	}
	if input.CounterpartyAccountID != nil {
		setter.CounterpartyAccountID = omitnull.From(*input.CounterpartyAccountID)
	}
	if input.CategoryID != nil {
		setter.CategoryID = omitnull.From(*input.CategoryID)
	}
	if input.Description != nil {
		setter.Description = omitnull.From(*input.Description)
	}
	return setter
}
