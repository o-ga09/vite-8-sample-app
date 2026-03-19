package repository_test

import (
	"os"
	"testing"

	infradb "github.com/o-ga09/vite-8-sample-app/internal/infra/db"
	"github.com/stephenafamo/bob"
)

// Set testDB to enable DB integration tests.
// Leave nil to skip all DB-dependent tests.
var testDB bob.Transactor[bob.Tx]

func TestMain(m *testing.M) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn != "" {
		db, err := infradb.NewDB(dsn)
		if err == nil {
			testDB = db
		}
	}
	os.Exit(m.Run())
}
