package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

var (
	ErrNotFound            = errors.New("not found")
	ErrDuplicate           = errors.New("duplicate key")
	ErrNilValue            = errors.New("nil value")
	ErrInvalidID           = errors.New("invalid id")
	ErrInvalidValue        = errors.New("invalid value")
	ErrForeignKeyViolation = errors.New("foreign key violation")
)

func wrapDBError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}

	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case "23505":
			return ErrDuplicate
		case "23503":
			return ErrForeignKeyViolation
		default:
			return fmt.Errorf("postgres error [%s]: %w", pqErr.Code, err)
		}
	}

	return err
}
