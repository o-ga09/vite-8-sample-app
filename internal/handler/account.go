package handler

import (
	"context"
	"errors"

	"github.com/o-ga09/vite-8-sample-app/internal/infra/dbgen"
	enums "github.com/o-ga09/vite-8-sample-app/internal/infra/dbgen/dbenums"
	oas "github.com/o-ga09/vite-8-sample-app/internal/oas"
	"github.com/o-ga09/vite-8-sample-app/internal/service"
	"github.com/shopspring/decimal"
)

// CreateAccount implements oas.Handler.
func (h *Handler) CreateAccount(ctx context.Context, req *oas.CreateAccountRequest, params oas.CreateAccountParams) (oas.CreateAccountRes, error) {
	initialBalance, err := decimal.NewFromString(req.InitialBalance)
	if err != nil {
		return badRequest("invalid initialBalance: " + err.Error()), nil
	}
	a, err := h.accountSvc.Create(ctx, params.WsId, req.Name, enums.AccountType(req.AccountType), initialBalance)
	if err != nil {
		return badRequest(err.Error()), nil
	}
	return toOASAccount(a), nil
}

// GetAccount implements oas.Handler.
func (h *Handler) GetAccount(ctx context.Context, params oas.GetAccountParams) (oas.GetAccountRes, error) {
	a, err := h.accountSvc.Get(ctx, params.WsId, params.AccountId)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return notFound("account not found"), nil
		}
		return notFound(err.Error()), nil
	}
	return toOASAccount(a), nil
}

// ListAccounts implements oas.Handler.
func (h *Handler) ListAccounts(ctx context.Context, params oas.ListAccountsParams) ([]oas.Account, error) {
	as, err := h.accountSvc.List(ctx, params.WsId)
	if err != nil {
		return nil, err
	}
	result := make([]oas.Account, 0, len(as))
	for _, a := range as {
		result = append(result, *toOASAccount(a))
	}
	return result, nil
}

// UpdateAccount implements oas.Handler.
func (h *Handler) UpdateAccount(ctx context.Context, req *oas.UpdateAccountRequest, params oas.UpdateAccountParams) (oas.UpdateAccountRes, error) {
	a, err := h.accountSvc.Update(ctx, params.WsId, params.AccountId, req.Name, enums.AccountType(req.AccountType))
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return (*oas.UpdateAccountNotFound)(notFound("account not found")), nil
		}
		return (*oas.UpdateAccountBadRequest)(badRequest(err.Error())), nil
	}
	return toOASAccount(a), nil
}

// DeleteAccount implements oas.Handler.
func (h *Handler) DeleteAccount(ctx context.Context, params oas.DeleteAccountParams) (oas.DeleteAccountRes, error) {
	err := h.accountSvc.Delete(ctx, params.WsId, params.AccountId)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return notFound("account not found"), nil
		}
		return notFound(err.Error()), nil
	}
	return &oas.DeleteAccountNoContent{}, nil
}

func toOASAccount(a *dbgen.Account) *oas.Account {
	return &oas.Account{
		ID:             a.ID,
		WorkspaceId:    a.WorkspaceID,
		Name:           a.Name,
		AccountType:    oas.AccountType(a.AccountType),
		InitialBalance: a.InitialBalance.String(),
	}
}
