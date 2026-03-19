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
	"github.com/stephenafamo/bob"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCategoryService_Create(t *testing.T) {
	wsID := uuid.MustParse("01960000-0000-7000-8000-400000000001")
	tests := []struct {
		name     string
		parentID *uuid.UUID
		mock     func() repository.CategoryRepository
		wantErr  bool
	}{
		{
			name:     "success: create root category",
			parentID: nil,
			mock: func() repository.CategoryRepository {
				return &CategoryRepositoryMock{
					CreateFunc: func(
						ctx context.Context,
						exec bob.Executor,
						setter *dbgen.CategorySetter,
					) (*dbgen.Category, error) {
						return &dbgen.Category{
							ID:          setter.ID.GetOrZero(),
							WorkspaceID: setter.WorkspaceID.GetOrZero(),
							Name:        setter.Name.GetOrZero(),
						}, nil
					},
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewCategoryService(bob.DB{}, tt.mock())
			got, err := svc.Create(context.Background(), wsID, "Food", enums.CategoryTypeExpense, tt.parentID)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, "Food", got.Name)
		})
	}
}

func TestCategoryService_Get_NotFound(t *testing.T) {
	wsID := uuid.MustParse("01960000-0000-7000-8000-400000000001")
	catID := uuid.MustParse("01960000-0000-7000-8000-400000000002")
	mock := &CategoryRepositoryMock{
		GetFunc: func(
			ctx context.Context,
			exec bob.Executor,
			workspaceID, id uuid.UUID,
		) (*dbgen.Category, error) {
			return nil, sql.ErrNoRows
		},
	}
	svc := service.NewCategoryService(bob.DB{}, mock)
	_, err := svc.Get(context.Background(), wsID, catID)
	require.Error(t, err)
	assert.ErrorIs(t, err, service.ErrNotFound)
}
