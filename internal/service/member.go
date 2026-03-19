package service

import (
	"context"

	"github.com/aarondl/opt/omit"
	"github.com/google/uuid"
	"github.com/o-ga09/vite-8-sample-app/internal/infra/dbgen"
	enums "github.com/o-ga09/vite-8-sample-app/internal/infra/dbgen/dbenums"
	"github.com/o-ga09/vite-8-sample-app/internal/repository"
	"github.com/stephenafamo/bob"
)

// MemberService handles member business logic.
type MemberService interface {
	Create(ctx context.Context, workspaceID uuid.UUID, email, displayName string, role enums.MemberRole) (*dbgen.Member, error)
	Get(ctx context.Context, workspaceID, id uuid.UUID) (*dbgen.Member, error)
	List(ctx context.Context, workspaceID uuid.UUID) (dbgen.MemberSlice, error)
	Update(ctx context.Context, workspaceID, id uuid.UUID, displayName string, role enums.MemberRole) (*dbgen.Member, error)
	Delete(ctx context.Context, workspaceID, id uuid.UUID) error
}

type memberService struct {
	db   bob.DB
	repo repository.MemberRepository
}

// NewMemberService creates a new MemberService.
func NewMemberService(db bob.DB, repo repository.MemberRepository) MemberService {
	return &memberService{db: db, repo: repo}
}

func (s *memberService) Create(ctx context.Context, workspaceID uuid.UUID, email, displayName string, role enums.MemberRole) (*dbgen.Member, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	setter := &dbgen.MemberSetter{
		ID:          omit.From(id),
		WorkspaceID: omit.From(workspaceID),
		Email:       omit.From(email),
		DisplayName: omit.From(displayName),
		Role:        omit.From(role),
	}
	return s.repo.Create(ctx, s.db, setter)
}

func (s *memberService) Get(ctx context.Context, workspaceID, id uuid.UUID) (*dbgen.Member, error) {
	m, err := s.repo.Get(ctx, s.db, workspaceID, id)
	if err != nil {
		return nil, mapRepoErr(err)
	}
	return m, nil
}

func (s *memberService) List(ctx context.Context, workspaceID uuid.UUID) (dbgen.MemberSlice, error) {
	return s.repo.List(ctx, s.db, workspaceID)
}

func (s *memberService) Update(ctx context.Context, workspaceID, id uuid.UUID, displayName string, role enums.MemberRole) (*dbgen.Member, error) {
	setter := &dbgen.MemberSetter{
		DisplayName: omit.From(displayName),
		Role:        omit.From(role),
	}
	m, err := s.repo.Update(ctx, s.db, workspaceID, id, setter)
	if err != nil {
		return nil, mapRepoErr(err)
	}
	return m, nil
}

func (s *memberService) Delete(ctx context.Context, workspaceID, id uuid.UUID) error {
	err := s.repo.Delete(ctx, s.db, workspaceID, id)
	if err != nil {
		return mapRepoErr(err)
	}
	return nil
}
