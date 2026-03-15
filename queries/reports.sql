-- name: GetCategoryExpenseSummary :many
SELECT
  c.id AS category_id,
  c.name AS category_name,
  COALESCE(SUM(t.amount), 0) AS total_expense
FROM categories c
LEFT JOIN transactions t
  ON t.workspace_id = c.workspace_id
  AND t.category_id = c.id
  AND t.transaction_type = 'expense'
  AND t.occurred_at >= $2
  AND t.occurred_at < $3
WHERE c.workspace_id = $1
GROUP BY c.id, c.name
ORDER BY total_expense DESC, c.name ASC;

-- name: GetAccountBalanceSummary :many
SELECT
  a.id AS account_id,
  a.name AS account_name,
  a.account_type,
  a.initial_balance
    + COALESCE(SUM(CASE
      WHEN t.transaction_type = 'income' AND t.account_id = a.id THEN t.amount
      WHEN t.transaction_type = 'expense' AND t.account_id = a.id THEN -t.amount
      WHEN t.transaction_type = 'transfer' AND t.counterparty_account_id = a.id THEN t.amount
      WHEN t.transaction_type = 'transfer' AND t.account_id = a.id THEN -t.amount
      ELSE 0
    END), 0) AS current_balance
FROM accounts a
LEFT JOIN transactions t
  ON t.workspace_id = a.workspace_id
  AND (t.account_id = a.id OR t.counterparty_account_id = a.id)
  AND t.occurred_at < $2
WHERE a.workspace_id = $1
GROUP BY a.id, a.name, a.account_type, a.initial_balance
ORDER BY a.name ASC;

-- name: GetWorkspaceDashboard :one
SELECT
  $1::uuid AS workspace_id,
  COALESCE(SUM(CASE WHEN t.transaction_type = 'income' THEN t.amount ELSE 0 END), 0) AS total_income,
  COALESCE(SUM(CASE WHEN t.transaction_type = 'expense' THEN t.amount ELSE 0 END), 0) AS total_expense,
  COALESCE(SUM(CASE WHEN t.transaction_type = 'income' THEN t.amount ELSE -t.amount END), 0) AS net_flow,
  COUNT(*)::bigint AS transaction_count
FROM transactions t
WHERE t.workspace_id = $1
  AND t.occurred_at >= $2
  AND t.occurred_at < $3;
