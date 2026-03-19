package handler

import (
	"context"
	"errors"

	"github.com/aarondl/opt/null"
	"github.com/o-ga09/vite-8-sample-app/internal/infra/dbgen"
	enums "github.com/o-ga09/vite-8-sample-app/internal/infra/dbgen/dbenums"
	oas "github.com/o-ga09/vite-8-sample-app/internal/oas"
	"github.com/o-ga09/vite-8-sample-app/internal/service"
	"github.com/shopspring/decimal"
)

// CreateTransaction implements oas.Handler.
func (h *Handler) CreateTransaction(ctx context.Context, req *oas.CreateTransactionRequest, params oas.CreateTransactionParams) (oas.CreateTransactionRes, error) {
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return badRequest("invalid amount: " + err.Error()), nil
	}

	input := service.CreateTransactionInput{
		TransactionType: enums.TransactionType(req.TransactionType),
		Amount:          amount,
		OccurredAt:      req.OccurredAt,
	}
	if v, ok := req.AccountId.Get(); ok {
		input.AccountID = &v
	}
	if v, ok := req.CounterpartyAccountId.Get(); ok {
		input.CounterpartyAccountID = &v
	}
	if v, ok := req.CategoryId.Get(); ok {
		input.CategoryID = &v
	}
	if v, ok := req.Description.Get(); ok {
		input.Description = &v
	}

	tx, err := h.transactionSvc.Create(ctx, params.WsId, input)
	if err != nil {
		return badRequest(err.Error()), nil
	}
	return toOASTransaction(tx), nil
}

// GetTransaction implements oas.Handler.
func (h *Handler) GetTransaction(ctx context.Context, params oas.GetTransactionParams) (oas.GetTransactionRes, error) {
	tx, err := h.transactionSvc.Get(ctx, params.WsId, params.TxId)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return notFound("transaction not found"), nil
		}
		return notFound(err.Error()), nil
	}
	return toOASTransaction(tx), nil
}

// ListTransactions implements oas.Handler.
func (h *Handler) ListTransactions(ctx context.Context, params oas.ListTransactionsParams) ([]oas.Transaction, error) {
	var (
		txs dbgen.TransactionSlice
		err error
	)
	if params.From.IsSet() && params.To.IsSet() {
		txs, err = h.transactionSvc.ListByPeriod(ctx, params.WsId, params.From.Value, params.To.Value)
	} else {
		txs, err = h.transactionSvc.List(ctx, params.WsId)
	}
	if err != nil {
		return nil, err
	}
	result := make([]oas.Transaction, 0, len(txs))
	for _, tx := range txs {
		result = append(result, *toOASTransaction(tx))
	}
	return result, nil
}

// UpdateTransaction implements oas.Handler.
func (h *Handler) UpdateTransaction(ctx context.Context, req *oas.UpdateTransactionRequest, params oas.UpdateTransactionParams) (oas.UpdateTransactionRes, error) {
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return (*oas.UpdateTransactionBadRequest)(badRequest("invalid amount: " + err.Error())), nil
	}

	input := service.CreateTransactionInput{
		TransactionType: enums.TransactionType(req.TransactionType),
		Amount:          amount,
		OccurredAt:      req.OccurredAt,
	}
	if v, ok := req.AccountId.Get(); ok {
		input.AccountID = &v
	}
	if v, ok := req.CounterpartyAccountId.Get(); ok {
		input.CounterpartyAccountID = &v
	}
	if v, ok := req.CategoryId.Get(); ok {
		input.CategoryID = &v
	}
	if v, ok := req.Description.Get(); ok {
		input.Description = &v
	}

	tx, err := h.transactionSvc.Update(ctx, params.WsId, params.TxId, input)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return (*oas.UpdateTransactionNotFound)(notFound("transaction not found")), nil
		}
		return (*oas.UpdateTransactionBadRequest)(badRequest(err.Error())), nil
	}
	return toOASTransaction(tx), nil
}

// DeleteTransaction implements oas.Handler.
func (h *Handler) DeleteTransaction(ctx context.Context, params oas.DeleteTransactionParams) (oas.DeleteTransactionRes, error) {
	err := h.transactionSvc.Delete(ctx, params.WsId, params.TxId)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return notFound("transaction not found"), nil
		}
		return notFound(err.Error()), nil
	}
	return &oas.DeleteTransactionNoContent{}, nil
}

func toOASTransaction(tx *dbgen.Transaction) *oas.Transaction {
	result := &oas.Transaction{
		ID:              tx.ID,
		WorkspaceId:     tx.WorkspaceID,
		TransactionType: oas.TransactionType(tx.TransactionType),
		Amount:          tx.Amount.String(),
		OccurredAt:      tx.OccurredAt,
	}
	if v, ok := tx.AccountID.Get(); ok {
		result.AccountId = nullValToOptNilUUID(null.From(v))
	}
	if v, ok := tx.CounterpartyAccountID.Get(); ok {
		result.CounterpartyAccountId = nullValToOptNilUUID(null.From(v))
	}
	if v, ok := tx.CategoryID.Get(); ok {
		result.CategoryId = nullValToOptNilUUID(null.From(v))
	}
	if v, ok := tx.Description.Get(); ok {
		result.Description = nullValToOptNilString(null.From(v))
	}
	return result
}

func nullValToOptNilString(v null.Val[string]) oas.OptNilString {
	val, ok := v.Get()
	if !ok {
		return oas.OptNilString{Set: true, Null: true}
	}
	return oas.OptNilString{Value: val, Set: true}
}
