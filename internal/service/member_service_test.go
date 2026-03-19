package service_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/o-ga09/vite-8-sample-app/internal/infra/dbgen"
	enums "github.com/o-ga09/vite-8-sample-app/internal/infra/dbgen/dbenums"
	"github.com/o-ga09/vite-8-sample-app/internal/repository"
	"github.com/o-ga09/vite-8-sample-app/internal/service"
	"github.com/stephenafamo/bob"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemberService_Create(t *testing.T) {
	wsID := uuid.MustParse("01960000-0000-7000-8000-200000000001")
	tests := []struct {
		name    string
		mock    func() repository.MemberRepository
		wantErr bool
	}{
		{
			name: "success: create member",
			mock: func() repository.MemberRepository {
				return &MemberRepositoryMock{
					CreateFunc: func(
						ctx context.Context,
						exec bob.Executor,
						setter *dbgen.MemberSetter,
					) (*dbgen.Member, error) {
						return &dbgen.Member{
							ID:          setter.ID.GetOrZero(),
							WorkspaceID: setter.WorkspaceID.GetOrZero(),
							Email:       setter.Email.GetOrZero(),
							DisplayName: setter.DisplayName.GetOrZero(),
							Role:        setter.Role.GetOrZero(),
						}, nil
					},
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewMemberService(bob.DB{}, tt.mock())
			got, err := svc.Create(context.Background(), wsID, "test@example.com", "Test User", enums.MemberRoleAdmin)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, "test@example.com", got.Email)
			assert.Equal(t, wsID, got.WorkspaceID)
		})
	}
}

func TestMemberService_Get_NotFound(t *testing.T) {
	wsID := uuid.MustParse("01960000-0000-7000-8000-200000000001")
	memberID := uuid.MustParse("01960000-0000-7000-8000-200000000002")
	mock := &MemberRepositoryMock{
		GetFunc: func(
			ctx context.Context,
			exec bob.Executor,
			workspaceID, id uuid.UUID,
		) (*dbgen.Member, error) {
			return nil, sql.ErrNoRows
		},
	}
	svc := service.NewMemberService(bob.DB{}, mock)
	_, err := svc.Get(context.Background(), wsID, memberID)
	require.Error(t, err)
	assert.ErrorIs(t, err, service.ErrNotFound)
}

func TestMemberService_Delete_NotFound(t *testing.T) {
	wsID := uuid.MustParse("01960000-0000-7000-8000-200000000001")
	memberID := uuid.MustParse("01960000-0000-7000-8000-200000000003")
	mock := &MemberRepositoryMock{
		DeleteFunc: func(
			ctx context.Context,
			exec bob.Executor,
			workspaceID, id uuid.UUID,
		) error {
			return sql.ErrNoRows
		},
	}
	svc := service.NewMemberService(bob.DB{}, mock)
	err := svc.Delete(context.Background(), wsID, memberID)
	require.Error(t, err)
	assert.ErrorIs(t, err, service.ErrNotFound)
}
