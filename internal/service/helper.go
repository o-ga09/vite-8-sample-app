package service

import (
	"database/sql"
	"errors"
)

// mapRepoErr maps repository errors to service-layer errors.
// sql.ErrNoRows is mapped to ErrNotFound.
func mapRepoErr(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	return err
}
