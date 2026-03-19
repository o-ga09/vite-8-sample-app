package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/o-ga09/vite-8-sample-app/internal/repository"
	"github.com/o-ga09/vite-8-sample-app/internal/service"
	"github.com/o-ga09/vite-8-sample-app/queries"
	"github.com/shopspring/decimal"
	"github.com/stephenafamo/bob"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReportService_GetCategoryExpenses(t *testing.T) {
	wsID := uuid.MustParse("01960000-0000-7000-8000-600000000001")
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		mock    func() repository.ReportRepository
		wantLen int
		wantErr bool
	}{
		{
			name: "success: get category expenses",
			mock: func() repository.ReportRepository {
				return &ReportRepositoryMock{
					GetCategoryExpenseSummaryFunc: func(
						ctx context.Context,
						exec bob.Executor,
						workspaceID uuid.UUID,
						from, to time.Time,
					) ([]queries.GetCategoryExpenseSummaryRow, error) {
						return []queries.GetCategoryExpenseSummaryRow{
							{CategoryName: "Food", TotalExpense: decimal.NewFromInt(1000)},
						}, nil
					},
				}
			},
			wantLen: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewReportService(bob.DB{}, tt.mock())
			got, err := svc.GetCategoryExpenses(context.Background(), wsID, from, to)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Len(t, got, tt.wantLen)
		})
	}
}
