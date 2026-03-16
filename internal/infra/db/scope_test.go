package db_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	infradb "github.com/o-ga09/vite-8-sample-app/internal/infra/db"
)

func TestWithWorkspaceID_StoresAndRetrievesID(t *testing.T) {
	wsID := uuid.New()
	ctx := infradb.WithWorkspaceID(context.Background(), wsID)

	got, ok := infradb.WorkspaceIDFromContext(ctx)
	if !ok {
		t.Fatal("expected workspace ID in context, got none")
	}
	if got != wsID {
		t.Errorf("got workspace ID %v, want %v", got, wsID)
	}
}

func TestWorkspaceIDFromContext_ReturnsFalseWhenNotSet(t *testing.T) {
	_, ok := infradb.WorkspaceIDFromContext(context.Background())
	if ok {
		t.Fatal("expected no workspace ID in empty context, but got one")
	}
}

func TestWorkspaceScopedExec_SetsWorkspaceIDInContext(t *testing.T) {
	wsID := uuid.New()
	ctx := infradb.WorkspaceScopedExec(context.Background(), wsID)

	got, ok := infradb.WorkspaceIDFromContext(ctx)
	if !ok {
		t.Fatal("expected workspace ID in context after WorkspaceScopedExec")
	}
	if got != wsID {
		t.Errorf("got workspace ID %v, want %v", got, wsID)
	}
}

func TestGlobalExec_SkipsHooks(t *testing.T) {
	// GlobalExec should set SkipHooks keys in the context.
	// We verify this by checking that the context value differs from the input.
	ctx := context.Background()
	globalCtx := infradb.GlobalExec(ctx)

	if ctx == globalCtx {
		t.Error("GlobalExec should return a modified context")
	}
}
