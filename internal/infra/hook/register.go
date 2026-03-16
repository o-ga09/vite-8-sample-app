package hook

import (
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql/dialect"
)

// RegisterHooks appends the workspace-scoped query hooks to the provided hook sets.
// Call this once per tenant-scoped bob table during application initialisation.
//
// Example (after bobgen models are generated):
//
//	hook.RegisterHooks(
//	    &models.Accounts.SelectQueryHooks,
//	    &models.Accounts.UpdateQueryHooks,
//	    &models.Accounts.DeleteQueryHooks,
//	)
func RegisterHooks(
	selectHooks *bob.Hooks[*dialect.SelectQuery, bob.SkipQueryHooksKey],
	updateHooks *bob.Hooks[*dialect.UpdateQuery, bob.SkipQueryHooksKey],
	deleteHooks *bob.Hooks[*dialect.DeleteQuery, bob.SkipQueryHooksKey],
) {
	if selectHooks != nil {
		selectHooks.AppendHooks(WorkspaceSelectHook)
	}
	if updateHooks != nil {
		updateHooks.AppendHooks(WorkspaceUpdateHook)
	}
	if deleteHooks != nil {
		deleteHooks.AppendHooks(WorkspaceDeleteHook)
	}
}
