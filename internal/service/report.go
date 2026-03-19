package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/o-ga09/vite-8-sample-app/internal/repository"
	"github.com/o-ga09/vite-8-sample-app/queries"
	"github.com/stephenafamo/bob"
)

// ReportService handles report business logic.
type ReportService interface {
	GetCategoryExpenses(ctx context.Context, workspaceID uuid.UUID, from, to time.Time) ([]queries.GetCategoryExpenseSummaryRow, error)
	GetAccountBalances(ctx context.Context, workspaceID uuid.UUID, asOf time.Time) ([]queries.GetAccountBalanceSummaryRow, error)
	GetDashboard(ctx context.Context, workspaceID uuid.UUID, from, to time.Time) ([]queries.GetWorkspaceDashboardRow, error)
}

type reportService struct {
	db   bob.DB
	repo repository.ReportRepository
}

// NewReportService creates a new ReportService.
func NewReportService(db bob.DB, repo repository.ReportRepository) ReportService {
	return &reportService{db: db, repo: repo}
}

func (s *reportService) GetCategoryExpenses(ctx context.Context, workspaceID uuid.UUID, from, to time.Time) ([]queries.GetCategoryExpenseSummaryRow, error) {
	return s.repo.GetCategoryExpenseSummary(ctx, s.db, workspaceID, from, to)
}

func (s *reportService) GetAccountBalances(ctx context.Context, workspaceID uuid.UUID, asOf time.Time) ([]queries.GetAccountBalanceSummaryRow, error) {
	return s.repo.GetAccountBalanceSummary(ctx, s.db, workspaceID, asOf)
}

func (s *reportService) GetDashboard(ctx context.Context, workspaceID uuid.UUID, from, to time.Time) ([]queries.GetWorkspaceDashboardRow, error) {
	return s.repo.GetWorkspaceDashboard(ctx, s.db, workspaceID, from, to)
}
