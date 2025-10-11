// Package pg implements Postgres-specific database support
// for the exceptions package.
package pg

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/WebDeveloperBen/ai-gateway/internal/exceptions"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// ErrorHandler wraps DB errors into domain-specific API errors.
type ErrorHandler func(error) error

// MakeErrorHandler returns an ErrorHandler customized for a resource.
func MakeErrorHandler(resource string) ErrorHandler {
	return func(err error) error {
		if err == nil {
			return nil
		}

		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return exceptions.New(http.StatusNotFound, "Not Found", fmt.Sprintf("%s not found", resource))

		case isUniqueViolation(err):
			return exceptions.New(http.StatusConflict, "Conflict", fmt.Sprintf("%s already exists", resource))

		case isForeignKeyViolation(err):
			return exceptions.New(http.StatusBadRequest, "Invalid Reference", fmt.Sprintf("invalid reference for %s", resource))

		case isCheckViolation(err):
			return exceptions.New(http.StatusBadRequest, "Check Violation", fmt.Sprintf("invalid value for %s", resource))

		case isNotNullViolation(err):
			return exceptions.New(http.StatusBadRequest, "Missing Required Field", fmt.Sprintf("missing required field for %s", resource))

		default:
			return err // fallback: bubble the DB error unchanged
		}
	}
}

// PGError helpers
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505" // unique_violation
}

func isForeignKeyViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23503" // foreign_key_violation
}

func isCheckViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23514" // check_violation
}

func isNotNullViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23502" // not_null_violation
}
