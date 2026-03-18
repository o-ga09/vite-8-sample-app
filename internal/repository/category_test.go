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
	"github.com/stephenafamo/bob"
)

// createCategoryDirect inserts a category without parent_id (root category)
// The factory for Category has a self-referential bug causing infinite recursion,
// so we insert directly using dbgen.
func createCategoryDirect(ctx context.Context, t *testing.T, exec bob.Executor, wsID uuid.UUID) *dbgen.Category {
	t.Helper()
	setter := &dbgen.CategorySetter{
		ID:           omit.From(uuid.New()),
		WorkspaceID:  omit.From(wsID),
		Name:         omit.From("カテゴリ-" + uuid.New().String()[:8]),
		CategoryType: omit.From(dbenums.CategoryTypeExpense),
	}
	cat, err := dbgen.Categories.Insert(setter).One(ctx, exec)
	if err != nil {
		t.Fatalf("createCategoryDirect: %v", err)
	}
	return cat
}

func TestNewCategoryRepository_ReturnsNonNil(t *testing.T) {
	r := repository.NewCategoryRepository()
	if r == nil {
		t.Fatal("expected non-nil CategoryRepository")
	}
}

func TestCategoryRepository_Create(t *testing.T) {
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

	r := repository.NewCategoryRepository()
	setter := &dbgen.CategorySetter{
		ID:           omit.From(uuid.New()),
		WorkspaceID:  omit.From(ws.ID),
		Name:         omit.From("食費"),
		CategoryType: omit.From(dbenums.CategoryTypeExpense),
	}

	cat, err := r.Create(ctx, tx, setter)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if cat == nil {
		t.Fatal("expected non-nil category")
	}
	if cat.WorkspaceID != ws.ID {
		t.Errorf("got WorkspaceID %v, want %v", cat.WorkspaceID, ws.ID)
	}
}

func TestCategoryRepository_Get_ExistingCategory(t *testing.T) {
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
	cat := createCategoryDirect(ctx, t, tx, ws.ID)

	r := repository.NewCategoryRepository()
	got, err := r.Get(ctx, tx, cat.WorkspaceID, cat.ID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got.ID != cat.ID {
		t.Errorf("got ID %v, want %v", got.ID, cat.ID)
	}
}

func TestCategoryRepository_Get_NonExistentID(t *testing.T) {
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

	r := repository.NewCategoryRepository()
	_, err = r.Get(ctx, tx, uuid.New(), uuid.New())
	if err == nil {
		t.Fatal("expected error for non-existent category, got nil")
	}
}

func TestCategoryRepository_List(t *testing.T) {
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

	for i := 0; i < 2; i++ {
		createCategoryDirect(ctx, t, tx, ws.ID)
	}

	r := repository.NewCategoryRepository()
	list, err := r.List(ctx, tx, ws.ID)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(list) < 2 {
		t.Errorf("expected at least 2 categories, got %d", len(list))
	}
}

func TestCategoryRepository_Update(t *testing.T) {
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
	cat := createCategoryDirect(ctx, t, tx, ws.ID)

	r := repository.NewCategoryRepository()
	const newName = "updated-category"
	setter := &dbgen.CategorySetter{
		Name: omit.From(newName),
	}

	updated, err := r.Update(ctx, tx, cat.WorkspaceID, cat.ID, setter)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if updated.Name != newName {
		t.Errorf("got Name %q, want %q", updated.Name, newName)
	}
}

func TestCategoryRepository_Delete(t *testing.T) {
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
	cat := createCategoryDirect(ctx, t, tx, ws.ID)

	r := repository.NewCategoryRepository()
	if err = r.Delete(ctx, tx, cat.WorkspaceID, cat.ID); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = r.Get(ctx, tx, cat.WorkspaceID, cat.ID)
	if err == nil {
		t.Fatal("expected error after delete, got nil")
	}
}
