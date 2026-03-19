package service_test

import (
	"context"
	"database/sql"
	"testing"

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

func TestAccountService_Create(t *testing.T) {
	wsID := uuid.MustParse("01960000-0000-7000-8000-300000000001")
	tests := []struct {
		name    string
		mock    func() repository.AccountRepository
		wantErr bool
	}{
		{
			name: "success: create account",
			mock: func() repository.AccountRepository {
				return &AccountRepositoryMock{
					CreateFunc: func(
						ctx context.Context,
						exec bob.Executor,
						setter *dbgen.AccountSetter,
					) (*dbgen.Account, error) {
						return &dbgen.Account{
							ID:             setter.ID.GetOrZero(),
							WorkspaceID:    setter.WorkspaceID.GetOrZero(),
							Name:           setter.Name.GetOrZero(),
							AccountType:    setter.AccountType.GetOrZero(),
							InitialBalance: setter.InitialBalance.GetOrZero(),
						}, nil
					},
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewAccountService(bob.DB{}, tt.mock())
			got, err := svc.Create(context.Background(), wsID, "Cash", enums.AccountTypeCash, decimal.NewFromInt(1000))
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, "Cash", got.Name)
			assert.Equal(t, wsID, got.WorkspaceID)
		})
	}
}

func TestAccountService_Get_NotFound(t *testing.T) {
	wsID := uuid.MustParse("01960000-0000-7000-8000-300000000001")
	accID := uuid.MustParse("01960000-0000-7000-8000-300000000002")
	mock := &AccountRepositoryMock{
		GetFunc: func(
			ctx context.Context,
			exec bob.Executor,
			workspaceID, id uuid.UUID,
		) (*dbgen.Account, error) {
			return nil, sql.ErrNoRows
		},
	}
	svc := service.NewAccountService(bob.DB{}, mock)
	_, err := svc.Get(context.Background(), wsID, accID)
	require.Error(t, err)
	assert.ErrorIs(t, err, service.ErrNotFound)
}
