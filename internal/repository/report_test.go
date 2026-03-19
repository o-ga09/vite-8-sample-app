package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/o-ga09/vite-8-sample-app/internal/infra/dbgen/factory"
	"github.com/o-ga09/vite-8-sample-app/internal/repository"
)

func TestNewReportRepository_ReturnsNonNil(t *testing.T) {
	r := repository.NewReportRepository()
	if r == nil {
		t.Fatal("expected non-nil ReportRepository")
	}
}

func TestReportRepository_GetCategoryExpenseSummary(t *testing.T) {
	if testDB == nil {
		t.Skip("skipping: no DB connection")
	}

	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)

	tx, err := testDB.Begin(ctx)
	if err != nil {
		t.Fatalf("failed to begin tx: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Insert a workspace to use as the scope.
	ws, err := factory.New().NewWorkspace().Create(ctx, tx)
	if err != nil {
		t.Fatalf("factory workspace: %v", err)
	}

	r := repository.NewReportRepository()
	now := time.Now()
	rows, err := r.GetCategoryExpenseSummary(ctx, tx, ws.ID, now.Add(-24*time.Hour), now)
	if err != nil {
		t.Fatalf("GetCategoryExpenseSummary failed: %v", err)
	}
	// Empty result is acceptable; we only verify the call succeeds.
	_ = rows
}

func TestReportRepository_GetAccountBalanceSummary(t *testing.T) {
	if testDB == nil {
		t.Skip("skipping: no DB connection")
	}

	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)

	tx, err := testDB.Begin(ctx)
	if err != nil {
		t.Fatalf("failed to begin tx: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	ws, err := factory.New().NewWorkspace().Create(ctx, tx)
	if err != nil {
		t.Fatalf("factory workspace: %v", err)
	}

	r := repository.NewReportRepository()
	rows, err := r.GetAccountBalanceSummary(ctx, tx, ws.ID, time.Now())
	if err != nil {
		t.Fatalf("GetAccountBalanceSummary failed: %v", err)
	}
	_ = rows
}

func TestReportRepository_GetWorkspaceDashboard(t *testing.T) {
	if testDB == nil {
		t.Skip("skipping: no DB connection")
	}

	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)

	tx, err := testDB.Begin(ctx)
	if err != nil {
		t.Fatalf("failed to begin tx: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	ws, err := factory.New().NewWorkspace().Create(ctx, tx)
	if err != nil {
		t.Fatalf("factory workspace: %v", err)
	}

	r := repository.NewReportRepository()
	now := time.Now()
	rows, err := r.GetWorkspaceDashboard(ctx, tx, ws.ID, now.Add(-30*24*time.Hour), now)
	if err != nil {
		t.Fatalf("GetWorkspaceDashboard failed: %v", err)
	}
	_ = rows
}

func TestReportRepository_GetCategoryExpenseSummary_UnknownWorkspace(t *testing.T) {
	if testDB == nil {
		t.Skip("skipping: no DB connection")
	}

	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)

	tx, err := testDB.Begin(ctx)
	if err != nil {
		t.Fatalf("failed to begin tx: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	r := repository.NewReportRepository()
	now := time.Now()
	// Non-existent workspace_id returns empty results, not an error.
	rows, err := r.GetCategoryExpenseSummary(ctx, tx, uuid.New(), now.Add(-24*time.Hour), now)
	if err != nil {
		t.Fatalf("unexpected error for unknown workspace: %v", err)
	}
	if len(rows) != 0 {
		t.Errorf("expected 0 rows for unknown workspace, got %d", len(rows))
	}
}
