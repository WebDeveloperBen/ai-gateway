package pg_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/WebDeveloperBen/ai-gateway/internal/exceptions"
	"github.com/WebDeveloperBen/ai-gateway/internal/exceptions/pg"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
)

func TestMakeErrorHandler(t *testing.T) {
	handler := pg.MakeErrorHandler("user")

	t.Run("nil error", func(t *testing.T) {
		result := handler(nil)
		require.NoError(t, result)
	})

	t.Run("ErrNoRows", func(t *testing.T) {
		result := handler(pgx.ErrNoRows)
		require.Error(t, result)

		apiErr, ok := result.(exceptions.APIError)
		require.True(t, ok)
		require.Equal(t, http.StatusNotFound, apiErr.Status())
		require.Equal(t, "Not Found", apiErr.Title())
		require.Equal(t, "user not found", apiErr.Detail())
	})

	t.Run("unique violation", func(t *testing.T) {
		pgErr := &pgconn.PgError{Code: "23505"}
		result := handler(pgErr)
		require.Error(t, result)

		apiErr, ok := result.(exceptions.APIError)
		require.True(t, ok)
		require.Equal(t, http.StatusConflict, apiErr.Status())
		require.Equal(t, "Conflict", apiErr.Title())
		require.Equal(t, "user already exists", apiErr.Detail())
	})

	t.Run("foreign key violation", func(t *testing.T) {
		pgErr := &pgconn.PgError{Code: "23503"}
		result := handler(pgErr)
		require.Error(t, result)

		apiErr, ok := result.(exceptions.APIError)
		require.True(t, ok)
		require.Equal(t, http.StatusBadRequest, apiErr.Status())
		require.Equal(t, "Invalid Reference", apiErr.Title())
		require.Equal(t, "invalid reference for user", apiErr.Detail())
	})

	t.Run("check violation", func(t *testing.T) {
		pgErr := &pgconn.PgError{Code: "23514"}
		result := handler(pgErr)
		require.Error(t, result)

		apiErr, ok := result.(exceptions.APIError)
		require.True(t, ok)
		require.Equal(t, http.StatusBadRequest, apiErr.Status())
		require.Equal(t, "Check Violation", apiErr.Title())
		require.Equal(t, "invalid value for user", apiErr.Detail())
	})

	t.Run("not null violation", func(t *testing.T) {
		pgErr := &pgconn.PgError{Code: "23502"}
		result := handler(pgErr)
		require.Error(t, result)

		apiErr, ok := result.(exceptions.APIError)
		require.True(t, ok)
		require.Equal(t, http.StatusBadRequest, apiErr.Status())
		require.Equal(t, "Missing Required Field", apiErr.Title())
		require.Equal(t, "missing required field for user", apiErr.Detail())
	})

	t.Run("unknown error", func(t *testing.T) {
		unknownErr := errors.New("unknown database error")
		result := handler(unknownErr)
		require.Equal(t, unknownErr, result)
	})

	t.Run("non-PgError with unique code", func(t *testing.T) {
		// Test that we don't panic on non-PgError types
		regularErr := errors.New("some error")
		result := handler(regularErr)
		require.Equal(t, regularErr, result)
	})
}

func TestErrorHandlerWithDifferentResources(t *testing.T) {
	tests := []struct {
		resource       string
		pgError        error
		expectedDetail string
	}{
		{"organization", pgx.ErrNoRows, "organization not found"},
		{"api_key", &pgconn.PgError{Code: "23505"}, "api_key already exists"},
		{"policy", &pgconn.PgError{Code: "23503"}, "invalid reference for policy"},
		{"model", &pgconn.PgError{Code: "23514"}, "invalid value for model"},
		{"user", &pgconn.PgError{Code: "23502"}, "missing required field for user"},
	}

	for _, tt := range tests {
		t.Run(tt.resource+"_error", func(t *testing.T) {
			handler := pg.MakeErrorHandler(tt.resource)
			result := handler(tt.pgError)

			apiErr, ok := result.(exceptions.APIError)
			require.True(t, ok)
			require.Equal(t, tt.expectedDetail, apiErr.Detail())
		})
	}
}
