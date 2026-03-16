-- migrate:up
CREATE TABLE accounts (
  id UUID PRIMARY KEY,
  workspace_id UUID NOT NULL,
  name TEXT NOT NULL,
  account_type account_type NOT NULL,
  initial_balance NUMERIC(14,2) NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(workspace_id, id),
  UNIQUE(workspace_id, name),
  CONSTRAINT fk_accounts_workspace
    FOREIGN KEY (workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE
);

-- migrate:down
DROP TABLE IF EXISTS accounts;
