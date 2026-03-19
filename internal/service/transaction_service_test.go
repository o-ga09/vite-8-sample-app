package service_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/o-ga09/vite-8-sample-app/internal/infra/dbgen"
	enums "github.com/o-ga09/vite-8-sample-app/internal/infra/dbgen/dbenums"
	"github.com/o-ga09/vite-8-sample-app/internal/repository"
	"github.com/o-ga09/vite-8-sample-app/internal/service"
	"github.com/shopspring/decimal"
	"github.com/stephenafamo/bob"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransactionService_Create(t *testing.T) {
	wsID := uuid.MustParse("01960000-0000-7000-8000-500000000001")
	accID := uuid.MustParse("01960000-0000-7000-8000-500000000002")
	tests := []struct {
		name    string
		mock    func() repository.TransactionRepository
		wantErr bool
	}{
		{
			name: "success: create transaction",
			mock: func() repository.TransactionRepository {
				return &TransactionRepositoryMock{
					CreateFunc: func(
						ctx context.Context,
						exec bob.Executor,
						setter *dbgen.TransactionSetter,
					) (*dbgen.Transaction, error) {
						return &dbgen.Transaction{
							ID:              setter.ID.GetOrZero(),
							WorkspaceID:     setter.WorkspaceID.GetOrZero(),
							TransactionType: setter.TransactionType.GetOrZero(),
							Amount:          setter.Amount.GetOrZero(),
						}, nil
					},
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewTransactionService(bob.DB{}, tt.mock())
			input := service.CreateTransactionInput{
				TransactionType: enums.TransactionTypeIncome,
				AccountID:       &accID,
				Amount:          decimal.NewFromInt(500),
				OccurredAt:      time.Now(),
			}
			got, err := svc.Create(context.Background(), wsID, input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, wsID, got.WorkspaceID)
			assert.Equal(t, decimal.NewFromInt(500), got.Amount)
		})
	}
}

func TestTransactionService_Get_NotFound(t *testing.T) {
	wsID := uuid.MustParse("01960000-0000-7000-8000-500000000001")
	txID := uuid.MustParse("01960000-0000-7000-8000-500000000003")
	mock := &TransactionRepositoryMock{
		GetFunc: func(
			ctx context.Context,
			exec bob.Executor,
			workspaceID, id uuid.UUID,
		) (*dbgen.Transaction, error) {
			return nil, sql.ErrNoRows
		},
	}
	svc := service.NewTransactionService(bob.DB{}, mock)
	_, err := svc.Get(context.Background(), wsID, txID)
	require.Error(t, err)
	assert.ErrorIs(t, err, service.ErrNotFound)
}
