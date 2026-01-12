package shared

import (
	"errors"
	"fmt"

	"github.com/lib/pq"
)

var (
	ErrUniqueViolation        = errors.New("unique constraint violation")
	ErrForeignKeyViolation    = errors.New("foreign key vioaltion")
	ErrNullConstrainViolation = errors.New("null constrain violation")
	ErrUserNotFound           = errors.New("Not Found")
)

func PostgresError(err error) error {
	var pqErr *pq.Error
	if !errors.As(err, &pqErr) {
		return err
	}

	switch pqErr.Code {
	case "23505":
		return fmt.Errorf("%w: %s", ErrUniqueViolation, pqErr.Detail)
	case "23503":
		return fmt.Errorf("%w: %s", ErrForeignKeyViolation, pqErr.Detail)
	case "23502":
		return fmt.Errorf("%w: %s", ErrNullConstrainViolation, pqErr.Detail)
	default:
		return fmt.Errorf("postgres err (%s): %w", pqErr.Code, err)
	}
}
