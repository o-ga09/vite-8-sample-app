package handler_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/o-ga09/vite-8-sample-app/internal/handler"
	"github.com/o-ga09/vite-8-sample-app/internal/infra/dbgen"
	oas "github.com/o-ga09/vite-8-sample-app/internal/oas"
	"github.com/o-ga09/vite-8-sample-app/internal/service"
)

type mockWorkspaceService struct {
	createFn func(ctx context.Context, name string) (*dbgen.Workspace, error)
	getFn    func(ctx context.Context, id uuid.UUID) (*dbgen.Workspace, error)
	listFn   func(ctx context.Context) (dbgen.WorkspaceSlice, error)
	updateFn func(ctx context.Context, id uuid.UUID, name string) (*dbgen.Workspace, error)
	deleteFn func(ctx context.Context, id uuid.UUID) error
}

func (m *mockWorkspaceService) Create(ctx context.Context, name string) (*dbgen.Workspace, error) {
	return m.createFn(ctx, name)
}
func (m *mockWorkspaceService) Get(ctx context.Context, id uuid.UUID) (*dbgen.Workspace, error) {
	return m.getFn(ctx, id)
}
func (m *mockWorkspaceService) List(ctx context.Context) (dbgen.WorkspaceSlice, error) {
	return m.listFn(ctx)
}
func (m *mockWorkspaceService) Update(ctx context.Context, id uuid.UUID, name string) (*dbgen.Workspace, error) {
	return m.updateFn(ctx, id, name)
}
func (m *mockWorkspaceService) Delete(ctx context.Context, id uuid.UUID) error {
	return m.deleteFn(ctx, id)
}

var _ service.WorkspaceService = (*mockWorkspaceService)(nil)

func TestHandler_CreateWorkspace(t *testing.T) {
	wsID := uuid.MustParse("01960000-0000-7000-8000-100000000001")
	mockSvc := &mockWorkspaceService{
		createFn: func(_ context.Context, name string) (*dbgen.Workspace, error) {
			return &dbgen.Workspace{ID: wsID, Name: name}, nil
		},
	}
	h := handler.New(handler.Options{WorkspaceService: mockSvc})
	req := &oas.CreateWorkspaceRequest{Name: "TestWS"}
	res, err := h.CreateWorkspace(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := res.(*oas.Workspace)
	if !ok {
		t.Fatalf("expected *oas.Workspace, got %T", res)
	}
	want := &oas.Workspace{ID: wsID, Name: "TestWS"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("CreateWorkspace mismatch (-want +got):\n%s", diff)
	}
}

func TestHandler_GetWorkspace_Found(t *testing.T) {
	wsID := uuid.MustParse("01960000-0000-7000-8000-100000000002")
	mockSvc := &mockWorkspaceService{
		getFn: func(_ context.Context, id uuid.UUID) (*dbgen.Workspace, error) {
			return &dbgen.Workspace{ID: id, Name: "found"}, nil
		},
	}
	h := handler.New(handler.Options{WorkspaceService: mockSvc})
	res, err := h.GetWorkspace(context.Background(), oas.GetWorkspaceParams{WsId: wsID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := res.(*oas.Workspace)
	if !ok {
		t.Fatalf("expected *oas.Workspace, got %T", res)
	}
	if got.ID != wsID {
		t.Errorf("expected ID %v, got %v", wsID, got.ID)
	}
}

func TestHandler_GetWorkspace_NotFound(t *testing.T) {
	wsID := uuid.MustParse("01960000-0000-7000-8000-100000000003")
	mockSvc := &mockWorkspaceService{
		getFn: func(_ context.Context, _ uuid.UUID) (*dbgen.Workspace, error) {
			return nil, service.ErrNotFound
		},
	}
	h := handler.New(handler.Options{WorkspaceService: mockSvc})
	res, err := h.GetWorkspace(context.Background(), oas.GetWorkspaceParams{WsId: wsID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.(*oas.ErrorResponse); !ok {
		t.Errorf("expected *oas.ErrorResponse for not found, got %T", res)
	}
}

func TestHandler_ListWorkspaces(t *testing.T) {
	mockSvc := &mockWorkspaceService{
		listFn: func(_ context.Context) (dbgen.WorkspaceSlice, error) {
			return dbgen.WorkspaceSlice{
				{ID: uuid.MustParse("01960000-0000-7000-8000-100000000004"), Name: "ws1"},
				{ID: uuid.MustParse("01960000-0000-7000-8000-100000000005"), Name: "ws2"},
			}, nil
		},
	}
	h := handler.New(handler.Options{WorkspaceService: mockSvc})
	res, err := h.ListWorkspaces(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 2 {
		t.Errorf("expected 2 workspaces, got %d", len(res))
	}
}

func TestHandler_DeleteWorkspace_OK(t *testing.T) {
	wsID := uuid.MustParse("01960000-0000-7000-8000-100000000006")
	mockSvc := &mockWorkspaceService{
		deleteFn: func(_ context.Context, _ uuid.UUID) error { return nil },
	}
	h := handler.New(handler.Options{WorkspaceService: mockSvc})
	res, err := h.DeleteWorkspace(context.Background(), oas.DeleteWorkspaceParams{WsId: wsID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.(*oas.DeleteWorkspaceNoContent); !ok {
		t.Errorf("expected *oas.DeleteWorkspaceNoContent, got %T", res)
	}
}

func TestHandler_DeleteWorkspace_NotFound(t *testing.T) {
	wsID := uuid.MustParse("01960000-0000-7000-8000-100000000007")
	mockSvc := &mockWorkspaceService{
		deleteFn: func(_ context.Context, _ uuid.UUID) error { return service.ErrNotFound },
	}
	h := handler.New(handler.Options{WorkspaceService: mockSvc})
	res, err := h.DeleteWorkspace(context.Background(), oas.DeleteWorkspaceParams{WsId: wsID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.(*oas.ErrorResponse); !ok {
		t.Errorf("expected *oas.ErrorResponse for not found, got %T", res)
	}
}
