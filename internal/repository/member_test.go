package repository_test

import (
	"context"
	"testing"

	"github.com/aarondl/opt/omit"
	"github.com/google/uuid"
	"github.com/o-ga09/vite-8-sample-app/internal/infra/dbgen"
	dbenums "github.com/o-ga09/vite-8-sample-app/internal/infra/dbgen/dbenums"
	"github.com/o-ga09/vite-8-sample-app/internal/infra/dbgen/factory"
	"github.com/o-ga09/vite-8-sample-app/internal/repository"
)

func TestNewMemberRepository_ReturnsNonNil(t *testing.T) {
	r := repository.NewMemberRepository()
	if r == nil {
		t.Fatal("expected non-nil MemberRepository")
	}
}

func TestMemberRepository_Create(t *testing.T) {
	if testDB == nil {
		t.Skip("skipping: no DB connection")
	}

	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)

	tx, err := testDB.Begin(ctx)
	if err != nil {
		t.Fatalf("failed to begin tx: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	ws, err := factory.New().NewWorkspace().Create(ctx, tx)
	if err != nil {
		t.Fatalf("factory workspace create failed: %v", err)
	}

	r := repository.NewMemberRepository()
	setter := &dbgen.MemberSetter{
		ID:          omit.From(uuid.New()),
		WorkspaceID: omit.From(ws.ID),
		Email:       omit.From("test@example.com"),
		DisplayName: omit.From("Test User"),
		Role:        omit.From(dbenums.MemberRoleMember),
	}

	member, err := r.Create(ctx, tx, setter)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if member == nil {
		t.Fatal("expected non-nil member")
	}
	if member.WorkspaceID != ws.ID {
		t.Errorf("got WorkspaceID %v, want %v", member.WorkspaceID, ws.ID)
	}
}

func TestMemberRepository_Get_ExistingMember(t *testing.T) {
	if testDB == nil {
		t.Skip("skipping: no DB connection")
	}

	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)

	tx, err := testDB.Begin(ctx)
	if err != nil {
		t.Fatalf("failed to begin tx: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	member, err := factory.New().NewMember().Create(ctx, tx)
	if err != nil {
		t.Fatalf("factory.Create failed: %v", err)
	}

	r := repository.NewMemberRepository()
	got, err := r.Get(ctx, tx, member.WorkspaceID, member.ID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got.ID != member.ID {
		t.Errorf("got ID %v, want %v", got.ID, member.ID)
	}
}

func TestMemberRepository_Get_NonExistentID(t *testing.T) {
	if testDB == nil {
		t.Skip("skipping: no DB connection")
	}

	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)

	tx, err := testDB.Begin(ctx)
	if err != nil {
		t.Fatalf("failed to begin tx: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	r := repository.NewMemberRepository()
	_, err = r.Get(ctx, tx, uuid.New(), uuid.New())
	if err == nil {
		t.Fatal("expected error for non-existent member, got nil")
	}
}

func TestMemberRepository_List(t *testing.T) {
	if testDB == nil {
		t.Skip("skipping: no DB connection")
	}

	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)

	tx, err := testDB.Begin(ctx)
	if err != nil {
		t.Fatalf("failed to begin tx: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	ws, err := factory.New().NewWorkspace().Create(ctx, tx)
	if err != nil {
		t.Fatalf("factory workspace create failed: %v", err)
	}

	f := factory.New()
	for i := 0; i < 2; i++ {
		if _, createErr := f.NewMember(factory.MemberMods.WithExistingWorkspace(ws)).Create(ctx, tx); createErr != nil {
			t.Fatalf("factory member create(%d) failed: %v", i, createErr)
		}
	}

	r := repository.NewMemberRepository()
	list, err := r.List(ctx, tx, ws.ID)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(list) < 2 {
		t.Errorf("expected at least 2 members, got %d", len(list))
	}
}

func TestMemberRepository_Update(t *testing.T) {
	if testDB == nil {
		t.Skip("skipping: no DB connection")
	}

	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)

	tx, err := testDB.Begin(ctx)
	if err != nil {
		t.Fatalf("failed to begin tx: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	member, err := factory.New().NewMember().Create(ctx, tx)
	if err != nil {
		t.Fatalf("factory.Create failed: %v", err)
	}

	r := repository.NewMemberRepository()
	const newDisplayName = "Updated User"
	setter := &dbgen.MemberSetter{
		DisplayName: omit.From(newDisplayName),
	}

	updated, err := r.Update(ctx, tx, member.WorkspaceID, member.ID, setter)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if updated.DisplayName != newDisplayName {
		t.Errorf("got DisplayName %q, want %q", updated.DisplayName, newDisplayName)
	}
}

func TestMemberRepository_Delete(t *testing.T) {
	if testDB == nil {
		t.Skip("skipping: no DB connection")
	}

	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)

	tx, err := testDB.Begin(ctx)
	if err != nil {
		t.Fatalf("failed to begin tx: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	member, err := factory.New().NewMember().Create(ctx, tx)
	if err != nil {
		t.Fatalf("factory.Create failed: %v", err)
	}

	r := repository.NewMemberRepository()
	if err = r.Delete(ctx, tx, member.WorkspaceID, member.ID); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = r.Get(ctx, tx, member.WorkspaceID, member.ID)
	if err == nil {
		t.Fatal("expected error after delete, got nil")
	}
}
