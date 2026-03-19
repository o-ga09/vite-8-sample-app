package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/o-ga09/vite-8-sample-app/queries"
	"github.com/stephenafamo/bob"
)

type ReportRepository interface {
	GetCategoryExpenseSummary(ctx context.Context, exec bob.Executor, workspaceID uuid.UUID, from, to time.Time) ([]queries.GetCategoryExpenseSummaryRow, error)
	GetAccountBalanceSummary(ctx context.Context, exec bob.Executor, workspaceID uuid.UUID, asOf time.Time) ([]queries.GetAccountBalanceSummaryRow, error)
	GetWorkspaceDashboard(ctx context.Context, exec bob.Executor, workspaceID uuid.UUID, from, to time.Time) ([]queries.GetWorkspaceDashboardRow, error)
}

type reportRepository struct{}

func NewReportRepository() ReportRepository {
	return &reportRepository{}
}

func (r *reportRepository) GetCategoryExpenseSummary(ctx context.Context, exec bob.Executor, workspaceID uuid.UUID, from, to time.Time) ([]queries.GetCategoryExpenseSummaryRow, error) {
	return queries.GetCategoryExpenseSummary(from, to, workspaceID).All(ctx, exec)
}

func (r *reportRepository) GetAccountBalanceSummary(ctx context.Context, exec bob.Executor, workspaceID uuid.UUID, asOf time.Time) ([]queries.GetAccountBalanceSummaryRow, error) {
	return queries.GetAccountBalanceSummary(asOf, workspaceID).All(ctx, exec)
}

func (r *reportRepository) GetWorkspaceDashboard(ctx context.Context, exec bob.Executor, workspaceID uuid.UUID, from, to time.Time) ([]queries.GetWorkspaceDashboardRow, error) {
	return queries.GetWorkspaceDashboard(workspaceID, from, to).All(ctx, exec)
}
