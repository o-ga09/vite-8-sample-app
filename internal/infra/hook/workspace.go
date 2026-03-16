package hook

import (
	"context"

	infradb "github.com/o-ga09/vite-8-sample-app/internal/infra/db"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dialect"
)

// WorkspaceSelectHook injects a WHERE workspace_id = $1 condition into SELECT queries.
// It reads the workspace_id from context using WorkspaceIDFromContext.
// If no workspace_id is set in the context, the hook is a no-op.
func WorkspaceSelectHook(ctx context.Context, _ bob.Executor, q *dialect.SelectQuery) (context.Context, error) {
	wsID, ok := infradb.WorkspaceIDFromContext(ctx)
	if !ok {
		return ctx, nil
	}
	q.AppendWhere(psql.Quote("workspace_id").EQ(psql.Arg(wsID)))
	return ctx, nil
}

// WorkspaceUpdateHook injects a WHERE workspace_id = $1 condition into UPDATE queries.
func WorkspaceUpdateHook(ctx context.Context, _ bob.Executor, q *dialect.UpdateQuery) (context.Context, error) {
	wsID, ok := infradb.WorkspaceIDFromContext(ctx)
	if !ok {
		return ctx, nil
	}
	q.AppendWhere(psql.Quote("workspace_id").EQ(psql.Arg(wsID)))
	return ctx, nil
}

// WorkspaceDeleteHook injects a WHERE workspace_id = $1 condition into DELETE queries.
func WorkspaceDeleteHook(ctx context.Context, _ bob.Executor, q *dialect.DeleteQuery) (context.Context, error) {
	wsID, ok := infradb.WorkspaceIDFromContext(ctx)
	if !ok {
		return ctx, nil
	}
	q.AppendWhere(psql.Quote("workspace_id").EQ(psql.Arg(wsID)))
	return ctx, nil
}
