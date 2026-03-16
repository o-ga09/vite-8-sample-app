package db_test

import (
	"testing"

	infradb "github.com/o-ga09/vite-8-sample-app/internal/infra/db"
)

func TestNewDB_ReturnsErrorOnInvalidDSN(t *testing.T) {
	_, err := infradb.NewDB("invalid-dsn")
	if err == nil {
		t.Fatal("expected error for invalid DSN, got nil")
	}
}

func TestNewDB_ReturnsErrorOnUnreachableHost(t *testing.T) {
	_, err := infradb.NewDB("postgres://user:pass@localhost:59999/nonexistent?sslmode=disable")
	if err == nil {
		t.Fatal("expected error for unreachable host, got nil")
	}
}
