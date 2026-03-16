package hook_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	infradb "github.com/o-ga09/vite-8-sample-app/internal/infra/db"
	"github.com/o-ga09/vite-8-sample-app/internal/infra/hook"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql/dialect"
)

func TestWorkspaceSelectHook_AppendsWhereCondition(t *testing.T) {
	wsID := uuid.New()
	ctx := infradb.WithWorkspaceID(context.Background(), wsID)

	q := &dialect.SelectQuery{}
	_, err := hook.WorkspaceSelectHook(ctx, nil, q)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(q.Where.Conditions) == 0 { //nolint:staticcheck
		t.Error("expected WHERE condition to be appended to SelectQuery")
	}
}

func TestWorkspaceSelectHook_NoOpWhenNoWorkspaceID(t *testing.T) {
	ctx := context.Background()

	q := &dialect.SelectQuery{}
	_, err := hook.WorkspaceSelectHook(ctx, nil, q)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(q.Where.Conditions) != 0 { //nolint:staticcheck
		t.Error("expected no WHERE condition when workspace ID is not set")
	}
}

func TestWorkspaceUpdateHook_AppendsWhereCondition(t *testing.T) {
	wsID := uuid.New()
	ctx := infradb.WithWorkspaceID(context.Background(), wsID)

	q := &dialect.UpdateQuery{}
	_, err := hook.WorkspaceUpdateHook(ctx, nil, q)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(q.Where.Conditions) == 0 { //nolint:staticcheck
		t.Error("expected WHERE condition to be appended to UpdateQuery")
	}
}

func TestWorkspaceDeleteHook_AppendsWhereCondition(t *testing.T) {
	wsID := uuid.New()
	ctx := infradb.WithWorkspaceID(context.Background(), wsID)

	q := &dialect.DeleteQuery{}
	_, err := hook.WorkspaceDeleteHook(ctx, nil, q)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(q.Where.Conditions) == 0 { //nolint:staticcheck
		t.Error("expected WHERE condition to be appended to DeleteQuery")
	}
}

func TestRegisterHooks_AppendsHooksToAllSets(t *testing.T) {
	var selectHooks bob.Hooks[*dialect.SelectQuery, bob.SkipQueryHooksKey]
	var updateHooks bob.Hooks[*dialect.UpdateQuery, bob.SkipQueryHooksKey]
	var deleteHooks bob.Hooks[*dialect.DeleteQuery, bob.SkipQueryHooksKey]

	hook.RegisterHooks(&selectHooks, &updateHooks, &deleteHooks)

	if len(selectHooks.GetHooks()) == 0 {
		t.Error("expected SelectQueryHooks to have hooks registered")
	}
	if len(updateHooks.GetHooks()) == 0 {
		t.Error("expected UpdateQueryHooks to have hooks registered")
	}
	if len(deleteHooks.GetHooks()) == 0 {
		t.Error("expected DeleteQueryHooks to have hooks registered")
	}
}
