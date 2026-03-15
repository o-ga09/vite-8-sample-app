-- migrate:up
CREATE TABLE categories (
  id UUID PRIMARY KEY,
  workspace_id UUID NOT NULL,
  name TEXT NOT NULL,
  category_type category_type NOT NULL,
  parent_id UUID,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(workspace_id, id),
  UNIQUE(workspace_id, category_type, name),
  CONSTRAINT fk_categories_workspace
    FOREIGN KEY (workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE,
  CONSTRAINT fk_categories_parent
    FOREIGN KEY (workspace_id, parent_id)
    REFERENCES categories(workspace_id, id)
    ON DELETE SET NULL
);

-- migrate:down
DROP TABLE IF EXISTS categories;
