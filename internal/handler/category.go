package handler

import (
	"context"
	"errors"

	"github.com/aarondl/opt/null"
	"github.com/google/uuid"
	"github.com/o-ga09/vite-8-sample-app/internal/infra/dbgen"
	enums "github.com/o-ga09/vite-8-sample-app/internal/infra/dbgen/dbenums"
	oas "github.com/o-ga09/vite-8-sample-app/internal/oas"
	"github.com/o-ga09/vite-8-sample-app/internal/service"
)

// CreateCategory implements oas.Handler.
func (h *Handler) CreateCategory(ctx context.Context, req *oas.CreateCategoryRequest, params oas.CreateCategoryParams) (oas.CreateCategoryRes, error) {
	var parentID *uuid.UUID
	if v, ok := req.ParentId.Get(); ok {
		parentID = &v
	}
	c, err := h.categorySvc.Create(ctx, params.WsId, req.Name, enums.CategoryType(req.CategoryType), parentID)
	if err != nil {
		return badRequest(err.Error()), nil
	}
	return toOASCategory(c), nil
}

// GetCategory implements oas.Handler.
func (h *Handler) GetCategory(ctx context.Context, params oas.GetCategoryParams) (oas.GetCategoryRes, error) {
	c, err := h.categorySvc.Get(ctx, params.WsId, params.CategoryId)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return notFound("category not found"), nil
		}
		return notFound(err.Error()), nil
	}
	return toOASCategory(c), nil
}

// ListCategories implements oas.Handler.
func (h *Handler) ListCategories(ctx context.Context, params oas.ListCategoriesParams) ([]oas.Category, error) {
	cs, err := h.categorySvc.List(ctx, params.WsId)
	if err != nil {
		return nil, err
	}
	result := make([]oas.Category, 0, len(cs))
	for _, c := range cs {
		result = append(result, *toOASCategory(c))
	}
	return result, nil
}

// UpdateCategory implements oas.Handler.
func (h *Handler) UpdateCategory(ctx context.Context, req *oas.UpdateCategoryRequest, params oas.UpdateCategoryParams) (oas.UpdateCategoryRes, error) {
	var parentID *uuid.UUID
	if v, ok := req.ParentId.Get(); ok {
		parentID = &v
	}
	c, err := h.categorySvc.Update(ctx, params.WsId, params.CategoryId, req.Name, enums.CategoryType(req.CategoryType), parentID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return (*oas.UpdateCategoryNotFound)(notFound("category not found")), nil
		}
		return (*oas.UpdateCategoryBadRequest)(badRequest(err.Error())), nil
	}
	return toOASCategory(c), nil
}

// DeleteCategory implements oas.Handler.
func (h *Handler) DeleteCategory(ctx context.Context, params oas.DeleteCategoryParams) (oas.DeleteCategoryRes, error) {
	err := h.categorySvc.Delete(ctx, params.WsId, params.CategoryId)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return notFound("category not found"), nil
		}
		return notFound(err.Error()), nil
	}
	return &oas.DeleteCategoryNoContent{}, nil
}

func toOASCategory(c *dbgen.Category) *oas.Category {
	result := &oas.Category{
		ID:           c.ID,
		WorkspaceId:  c.WorkspaceID,
		Name:         c.Name,
		CategoryType: oas.CategoryType(c.CategoryType),
	}
	if v, ok := c.ParentID.Get(); ok {
		result.ParentId = nullValToOptNilUUID(null.From(v))
	}
	return result
}

func nullValToOptNilUUID(v null.Val[uuid.UUID]) oas.OptNilUUID {
	val, ok := v.Get()
	if !ok {
		return oas.OptNilUUID{Set: true, Null: true}
	}
	return oas.OptNilUUID{Value: val, Set: true}
}
