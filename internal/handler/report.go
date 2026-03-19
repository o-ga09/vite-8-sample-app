package handler

import (
	"context"

	oas "github.com/o-ga09/vite-8-sample-app/internal/oas"
)

// GetDashboard implements oas.Handler.
func (h *Handler) GetDashboard(ctx context.Context, params oas.GetDashboardParams) (*oas.DashboardReport, error) {
	rows, err := h.reportSvc.GetDashboard(ctx, params.WsId, params.From, params.To)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return &oas.DashboardReport{
			WorkspaceId: params.WsId,
		}, nil
	}
	row := rows[0]
	return &oas.DashboardReport{
		WorkspaceId:      row.WorkspaceID,
		TotalIncome:      row.TotalIncome.String(),
		TotalExpense:     row.TotalExpense.String(),
		NetFlow:          row.NetFlow.String(),
		TransactionCount: row.TransactionCount,
	}, nil
}

// GetCategoryExpenses implements oas.Handler.
func (h *Handler) GetCategoryExpenses(ctx context.Context, params oas.GetCategoryExpensesParams) ([]oas.CategoryExpenseItem, error) {
	rows, err := h.reportSvc.GetCategoryExpenses(ctx, params.WsId, params.From, params.To)
	if err != nil {
		return nil, err
	}
	result := make([]oas.CategoryExpenseItem, 0, len(rows))
	for _, row := range rows {
		result = append(result, oas.CategoryExpenseItem{
			CategoryId:   row.CategoryID,
			CategoryName: row.CategoryName,
			TotalExpense: row.TotalExpense.String(),
		})
	}
	return result, nil
}

// GetAccountBalances implements oas.Handler.
func (h *Handler) GetAccountBalances(ctx context.Context, params oas.GetAccountBalancesParams) ([]oas.AccountBalanceItem, error) {
	rows, err := h.reportSvc.GetAccountBalances(ctx, params.WsId, params.AsOf)
	if err != nil {
		return nil, err
	}
	result := make([]oas.AccountBalanceItem, 0, len(rows))
	for _, row := range rows {
		result = append(result, oas.AccountBalanceItem{
			AccountId:      row.AccountID,
			AccountName:    row.AccountName,
			AccountType:    oas.AccountType(row.AccountType),
			CurrentBalance: row.CurrentBalance.String(),
		})
	}
	return result, nil
}
