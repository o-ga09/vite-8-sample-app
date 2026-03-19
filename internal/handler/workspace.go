package handler

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/o-ga09/vite-8-sample-app/internal/infra/dbgen"
	oas "github.com/o-ga09/vite-8-sample-app/internal/oas"
	"github.com/o-ga09/vite-8-sample-app/internal/service"
)

// CreateWorkspace implements oas.Handler.
func (h *Handler) CreateWorkspace(ctx context.Context, req *oas.CreateWorkspaceRequest) (oas.CreateWorkspaceRes, error) {
	ws, err := h.workspaceSvc.Create(ctx, req.Name)
	if err != nil {
		return badRequest(err.Error()), nil
	}
	return toOASWorkspace(ws), nil
}

// GetWorkspace implements oas.Handler.
func (h *Handler) GetWorkspace(ctx context.Context, params oas.GetWorkspaceParams) (oas.GetWorkspaceRes, error) {
	ws, err := h.workspaceSvc.Get(ctx, params.WsId)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return notFound("workspace not found"), nil
		}
		return notFound(err.Error()), nil
	}
	return toOASWorkspace(ws), nil
}

// ListWorkspaces implements oas.Handler.
func (h *Handler) ListWorkspaces(ctx context.Context) ([]oas.Workspace, error) {
	wss, err := h.workspaceSvc.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]oas.Workspace, 0, len(wss))
	for _, ws := range wss {
		result = append(result, *toOASWorkspace(ws))
	}
	return result, nil
}

// UpdateWorkspace implements oas.Handler.
func (h *Handler) UpdateWorkspace(ctx context.Context, req *oas.UpdateWorkspaceRequest, params oas.UpdateWorkspaceParams) (oas.UpdateWorkspaceRes, error) {
	ws, err := h.workspaceSvc.Update(ctx, params.WsId, req.Name)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return (*oas.UpdateWorkspaceNotFound)(notFound("workspace not found")), nil
		}
		return (*oas.UpdateWorkspaceBadRequest)(badRequest(err.Error())), nil
	}
	return toOASWorkspace(ws), nil
}

// DeleteWorkspace implements oas.Handler.
func (h *Handler) DeleteWorkspace(ctx context.Context, params oas.DeleteWorkspaceParams) (oas.DeleteWorkspaceRes, error) {
	err := h.workspaceSvc.Delete(ctx, params.WsId)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return notFound("workspace not found"), nil
		}
		return notFound(err.Error()), nil
	}
	return &oas.DeleteWorkspaceNoContent{}, nil
}

func toOASWorkspace(ws *dbgen.Workspace) *oas.Workspace {
	return &oas.Workspace{
		ID:   uuid.UUID(ws.ID),
		Name: ws.Name,
	}
}
