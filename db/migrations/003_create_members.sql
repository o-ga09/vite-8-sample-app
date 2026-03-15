-- migrate:up
CREATE TABLE members (
  id UUID PRIMARY KEY,
  workspace_id UUID NOT NULL,
  email TEXT NOT NULL,
  display_name TEXT NOT NULL,
  role member_role NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(workspace_id, id),
  UNIQUE(workspace_id, email),
  CONSTRAINT fk_members_workspace
    FOREIGN KEY (workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE
);

-- migrate:down
DROP TABLE IF EXISTS members;
