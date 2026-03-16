package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/stephenafamo/bob"
)

type workspaceIDCtxKey struct{}

// WithWorkspaceID stores workspace_id in context so QueryHooks can inject it.
func WithWorkspaceID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, workspaceIDCtxKey{}, id)
}

// WorkspaceIDFromContext retrieves the workspace_id from context.
func WorkspaceIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(workspaceIDCtxKey{}).(uuid.UUID)
	return id, ok
}

// WorkspaceScopedExec returns a context that includes the workspace_id for
// tenant-scoped QueryHooks to inject WHERE workspace_id = $1 automatically.
func WorkspaceScopedExec(ctx context.Context, wsID uuid.UUID) context.Context {
	return WithWorkspaceID(ctx, wsID)
}

// GlobalExec returns a context that bypasses QueryHooks.
// Use this for workspace-level operations (e.g., creating or listing workspaces)
// where tenant scoping should not apply.
func GlobalExec(ctx context.Context) context.Context {
	return bob.SkipHooks(ctx)
}
