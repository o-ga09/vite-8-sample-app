package service_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/o-ga09/vite-8-sample-app/internal/infra/dbgen"
	"github.com/o-ga09/vite-8-sample-app/internal/repository"
	"github.com/o-ga09/vite-8-sample-app/internal/service"
	"github.com/stephenafamo/bob"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkspaceService_Create(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		mock    func() repository.WorkspaceRepository
		wantErr bool
	}{
		{
			name:  "success: create workspace",
			input: "MyWorkspace",
			mock: func() repository.WorkspaceRepository {
				return &WorkspaceRepositoryMock{
					CreateFunc: func(
						ctx context.Context,
						exec bob.Executor,
						setter *dbgen.WorkspaceSetter,
					) (*dbgen.Workspace, error) {
						return &dbgen.Workspace{
							ID:   setter.ID.GetOrZero(),
							Name: setter.Name.GetOrZero(),
						}, nil
					},
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewWorkspaceService(bob.DB{}, tt.mock())
			got, err := svc.Create(context.Background(), tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.input, got.Name)
			assert.NotEqual(t, uuid.UUID{}, got.ID)
		})
	}
}

func TestWorkspaceService_Get(t *testing.T) {
	wsID := uuid.MustParse("01960000-0000-7000-8000-100000000001")
	tests := []struct {
		name    string
		id      uuid.UUID
		mock    func() repository.WorkspaceRepository
		wantErr bool
		wantID  uuid.UUID
	}{
		{
			name: "success: get by id",
			id:   wsID,
			mock: func() repository.WorkspaceRepository {
				return &WorkspaceRepositoryMock{
					GetFunc: func(
						ctx context.Context,
						exec bob.Executor,
						id uuid.UUID,
					) (*dbgen.Workspace, error) {
						return &dbgen.Workspace{ID: id, Name: "found"}, nil
					},
				}
			},
			wantID: wsID,
		},
		{
			name: "error: not found",
			id:   wsID,
			mock: func() repository.WorkspaceRepository {
				return &WorkspaceRepositoryMock{
					GetFunc: func(
						ctx context.Context,
						exec bob.Executor,
						id uuid.UUID,
					) (*dbgen.Workspace, error) {
						return nil, sql.ErrNoRows
					},
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewWorkspaceService(bob.DB{}, tt.mock())
			got, err := svc.Get(context.Background(), tt.id)
			if tt.wantErr {
				require.Error(t, err)
				assert.ErrorIs(t, err, service.ErrNotFound)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantID, got.ID)
		})
	}
}

func TestWorkspaceService_List(t *testing.T) {
	tests := []struct {
		name    string
		mock    func() repository.WorkspaceRepository
		wantLen int
	}{
		{
			name: "success: list workspaces",
			mock: func() repository.WorkspaceRepository {
				return &WorkspaceRepositoryMock{
					ListFunc: func(
						ctx context.Context,
						exec bob.Executor,
					) (dbgen.WorkspaceSlice, error) {
						return dbgen.WorkspaceSlice{
							{ID: uuid.New(), Name: "ws1"},
							{ID: uuid.New(), Name: "ws2"},
						}, nil
					},
				}
			},
			wantLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewWorkspaceService(bob.DB{}, tt.mock())
			got, err := svc.List(context.Background())
			require.NoError(t, err)
			assert.Len(t, got, tt.wantLen)
		})
	}
}

func TestWorkspaceService_Update(t *testing.T) {
	wsID := uuid.MustParse("01960000-0000-7000-8000-100000000002")
	tests := []struct {
		name     string
		id       uuid.UUID
		newName  string
		mock     func() repository.WorkspaceRepository
		wantErr  bool
		wantName string
	}{
		{
			name:    "success: update name",
			id:      wsID,
			newName: "UpdatedName",
			mock: func() repository.WorkspaceRepository {
				return &WorkspaceRepositoryMock{
					UpdateFunc: func(
						ctx context.Context,
						exec bob.Executor,
						id uuid.UUID,
						setter *dbgen.WorkspaceSetter,
					) (*dbgen.Workspace, error) {
						return &dbgen.Workspace{ID: id, Name: setter.Name.GetOrZero()}, nil
					},
				}
			},
			wantName: "UpdatedName",
		},
		{
			name:    "error: not found",
			id:      wsID,
			newName: "X",
			mock: func() repository.WorkspaceRepository {
				return &WorkspaceRepositoryMock{
					UpdateFunc: func(
						ctx context.Context,
						exec bob.Executor,
						id uuid.UUID,
						setter *dbgen.WorkspaceSetter,
					) (*dbgen.Workspace, error) {
						return nil, sql.ErrNoRows
					},
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewWorkspaceService(bob.DB{}, tt.mock())
			got, err := svc.Update(context.Background(), tt.id, tt.newName)
			if tt.wantErr {
				require.Error(t, err)
				assert.ErrorIs(t, err, service.ErrNotFound)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantName, got.Name)
		})
	}
}

func TestWorkspaceService_Delete(t *testing.T) {
	wsID := uuid.MustParse("01960000-0000-7000-8000-100000000003")
	tests := []struct {
		name    string
		id      uuid.UUID
		mock    func() repository.WorkspaceRepository
		wantErr bool
	}{
		{
			name: "success: delete workspace",
			id:   wsID,
			mock: func() repository.WorkspaceRepository {
				return &WorkspaceRepositoryMock{
					DeleteFunc: func(
						ctx context.Context,
						exec bob.Executor,
						id uuid.UUID,
					) error {
						return nil
					},
				}
			},
		},
		{
			name: "error: not found",
			id:   wsID,
			mock: func() repository.WorkspaceRepository {
				return &WorkspaceRepositoryMock{
					DeleteFunc: func(
						ctx context.Context,
						exec bob.Executor,
						id uuid.UUID,
					) error {
						return sql.ErrNoRows
					},
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewWorkspaceService(bob.DB{}, tt.mock())
			err := svc.Delete(context.Background(), tt.id)
			if tt.wantErr {
				require.Error(t, err)
				assert.ErrorIs(t, err, service.ErrNotFound)
				return
			}
			require.NoError(t, err)
		})
	}
}
