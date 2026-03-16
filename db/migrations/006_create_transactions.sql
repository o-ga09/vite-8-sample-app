-- migrate:up
CREATE TABLE transactions (
  id UUID PRIMARY KEY,
  workspace_id UUID NOT NULL,
  transaction_type transaction_type NOT NULL,
  account_id UUID,
  counterparty_account_id UUID,
  category_id UUID,
  amount NUMERIC(14,2) NOT NULL CHECK (amount > 0),
  occurred_at TIMESTAMPTZ NOT NULL,
  description TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(workspace_id, id),
  CONSTRAINT fk_transactions_workspace
    FOREIGN KEY (workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE,
  CONSTRAINT fk_transactions_account
    FOREIGN KEY (workspace_id, account_id)
    REFERENCES accounts(workspace_id, id)
    ON DELETE RESTRICT,
  CONSTRAINT fk_transactions_counterparty_account
    FOREIGN KEY (workspace_id, counterparty_account_id)
    REFERENCES accounts(workspace_id, id)
    ON DELETE RESTRICT,
  CONSTRAINT fk_transactions_category
    FOREIGN KEY (workspace_id, category_id)
    REFERENCES categories(workspace_id, id)
    ON DELETE RESTRICT,
  CONSTRAINT chk_transactions_type_relations CHECK (
    (
      transaction_type = 'transfer'
      AND account_id IS NOT NULL
      AND counterparty_account_id IS NOT NULL
      AND account_id <> counterparty_account_id
      AND category_id IS NULL
    )
    OR
    (
      transaction_type IN ('income', 'expense')
      AND account_id IS NOT NULL
      AND counterparty_account_id IS NULL
      AND category_id IS NOT NULL
    )
  )
);

-- migrate:down
DROP TABLE IF EXISTS transactions;
