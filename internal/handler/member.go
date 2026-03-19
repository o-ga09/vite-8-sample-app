package handler

import (
	"context"
	"errors"

	"github.com/o-ga09/vite-8-sample-app/internal/infra/dbgen"
	enums "github.com/o-ga09/vite-8-sample-app/internal/infra/dbgen/dbenums"
	oas "github.com/o-ga09/vite-8-sample-app/internal/oas"
	"github.com/o-ga09/vite-8-sample-app/internal/service"
)

// CreateMember implements oas.Handler.
func (h *Handler) CreateMember(ctx context.Context, req *oas.CreateMemberRequest, params oas.CreateMemberParams) (oas.CreateMemberRes, error) {
	m, err := h.memberSvc.Create(ctx, params.WsId, req.Email, req.DisplayName, enums.MemberRole(req.Role))
	if err != nil {
		return badRequest(err.Error()), nil
	}
	return toOASMember(m), nil
}

// GetMember implements oas.Handler.
func (h *Handler) GetMember(ctx context.Context, params oas.GetMemberParams) (oas.GetMemberRes, error) {
	m, err := h.memberSvc.Get(ctx, params.WsId, params.MemberId)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return notFound("member not found"), nil
		}
		return notFound(err.Error()), nil
	}
	return toOASMember(m), nil
}

// ListMembers implements oas.Handler.
func (h *Handler) ListMembers(ctx context.Context, params oas.ListMembersParams) ([]oas.Member, error) {
	ms, err := h.memberSvc.List(ctx, params.WsId)
	if err != nil {
		return nil, err
	}
	result := make([]oas.Member, 0, len(ms))
	for _, m := range ms {
		result = append(result, *toOASMember(m))
	}
	return result, nil
}

// UpdateMember implements oas.Handler.
func (h *Handler) UpdateMember(ctx context.Context, req *oas.UpdateMemberRequest, params oas.UpdateMemberParams) (oas.UpdateMemberRes, error) {
	m, err := h.memberSvc.Update(ctx, params.WsId, params.MemberId, req.DisplayName, enums.MemberRole(req.Role))
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return (*oas.UpdateMemberNotFound)(notFound("member not found")), nil
		}
		return (*oas.UpdateMemberBadRequest)(badRequest(err.Error())), nil
	}
	return toOASMember(m), nil
}

// DeleteMember implements oas.Handler.
func (h *Handler) DeleteMember(ctx context.Context, params oas.DeleteMemberParams) (oas.DeleteMemberRes, error) {
	err := h.memberSvc.Delete(ctx, params.WsId, params.MemberId)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return notFound("member not found"), nil
		}
		return notFound(err.Error()), nil
	}
	return &oas.DeleteMemberNoContent{}, nil
}

func toOASMember(m *dbgen.Member) *oas.Member {
	return &oas.Member{
		ID:          m.ID,
		WorkspaceId: m.WorkspaceID,
		Email:       m.Email,
		DisplayName: m.DisplayName,
		Role:        oas.MemberRole(m.Role),
	}
}
