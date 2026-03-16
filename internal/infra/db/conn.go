package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/stephenafamo/bob"
)

// NewDB opens and validates a PostgreSQL connection using the provided DSN.
// It wraps the *sql.DB with bob.DB to enable type-safe query building.
func NewDB(dsn string) (bob.DB, error) {
	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		return bob.DB{}, fmt.Errorf("open db: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		_ = sqlDB.Close()
		return bob.DB{}, fmt.Errorf("ping db: %w", err)
	}

	return bob.NewDB(sqlDB), nil
}
