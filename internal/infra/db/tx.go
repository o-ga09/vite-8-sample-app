package db

import "context"

// Tx is an interface for database transaction lifecycle.
// *sql.Tx satisfies this interface.
type Tx interface {
	Commit() error
	Rollback() error
}

// TxFunc is the function executed within a transaction body.
type TxFunc func(ctx context.Context) error

// WithTx executes fn within the provided transaction.
// If fn returns an error or panics, the transaction is rolled back.
// Otherwise, the transaction is committed.
//
// Callers are responsible for beginning the transaction before calling WithTx.
// The SQL/bob executor can be captured in fn's closure.
//
// Example:
//
//	sqlTx, err := db.DB.BeginTx(ctx, nil)
//	if err != nil { return err }
//	bobTx := bob.Tx{Tx: sqlTx}
//	return infradb.WithTx(ctx, sqlTx, func(ctx context.Context) error {
//	    return repo.DoSomething(ctx, bobTx)
//	})
func WithTx(ctx context.Context, tx Tx, fn TxFunc) error {
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(ctx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
