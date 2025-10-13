package exceptions_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/WebDeveloperBen/ai-gateway/internal/exceptions"
	"github.com/danielgtaylor/huma/v2"
	"github.com/stretchr/testify/require"
)

func TestAPIError(t *testing.T) {
	t.Run("NotFound", func(t *testing.T) {
		err := exceptions.NotFound("user not found")
		require.Equal(t, http.StatusNotFound, err.Status())
		require.Equal(t, "Not Found", err.Title())
		require.Equal(t, "user not found", err.Detail())
		require.Empty(t, err.ErrorDetails())
		require.Contains(t, err.Error(), "404 Not Found: user not found")
	})

	t.Run("Validation", func(t *testing.T) {
		fields := []huma.ErrorDetail{
			{Message: "invalid format", Location: "body.email"},
		}
		err := exceptions.Validation("validation failed", fields)
		require.Equal(t, http.StatusUnprocessableEntity, err.Status())
		require.Equal(t, "Validation Failed", err.Title())
		require.Equal(t, "validation failed", err.Detail())
		require.Len(t, err.ErrorDetails(), 1)
		require.Equal(t, "invalid format", err.ErrorDetails()[0].Message)
		require.Equal(t, "body.email", err.ErrorDetails()[0].Location)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		err := exceptions.Unauthorized("invalid credentials")
		require.Equal(t, http.StatusUnauthorized, err.Status())
		require.Equal(t, "Unauthorized", err.Title())
		require.Equal(t, "invalid credentials", err.Detail())
	})

	t.Run("Forbidden", func(t *testing.T) {
		err := exceptions.Forbidden("access denied")
		require.Equal(t, http.StatusForbidden, err.Status())
		require.Equal(t, "Forbidden", err.Title())
		require.Equal(t, "access denied", err.Detail())
	})

	t.Run("InternalServerError", func(t *testing.T) {
		err := exceptions.InternalServerError("database connection failed")
		require.Equal(t, http.StatusInternalServerError, err.Status())
		require.Equal(t, "Internal Server Error", err.Title())
		require.Equal(t, "database connection failed", err.Detail())
	})
}

func TestErrorHelpers(t *testing.T) {
	t.Run("IsNotFound", func(t *testing.T) {
		require.True(t, exceptions.IsNotFound(exceptions.NotFound("test")))
		require.False(t, exceptions.IsNotFound(exceptions.Unauthorized("test")))
		require.False(t, exceptions.IsNotFound(errors.New("regular error")))
	})

	t.Run("IsUnauthorized", func(t *testing.T) {
		require.True(t, exceptions.IsUnauthorized(exceptions.Unauthorized("test")))
		require.False(t, exceptions.IsUnauthorized(exceptions.NotFound("test")))
		require.False(t, exceptions.IsUnauthorized(errors.New("regular error")))
	})

	t.Run("IsValidationError", func(t *testing.T) {
		require.True(t, exceptions.IsValidationError(exceptions.Validation("test", nil)))
		require.False(t, exceptions.IsValidationError(exceptions.NotFound("test")))
		require.False(t, exceptions.IsValidationError(errors.New("regular error")))
	})

	t.Run("Maybe", func(t *testing.T) {
		// Test with matching condition - should return fallback
		result := exceptions.Maybe(
			exceptions.NotFound("original"),
			exceptions.IsNotFound,
			exceptions.InternalServerError("fallback"),
		)
		require.Equal(t, exceptions.InternalServerError("fallback"), result)

		// Test with non-matching condition - should return original
		result = exceptions.Maybe(
			exceptions.Unauthorized("original"),
			exceptions.IsNotFound,
			exceptions.InternalServerError("fallback"),
		)
		require.Equal(t, exceptions.Unauthorized("original"), result)
	})
}

func TestHandle(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		handler := exceptions.Handle(func(ctx context.Context, input *string) (*string, error) {
			result := "success"
			return &result, nil
		})

		result, err := handler(context.Background(), stringPtr("input"))
		require.NoError(t, err)
		require.Equal(t, "success", *result)
	})

	t.Run("APIError conversion", func(t *testing.T) {
		handler := exceptions.Handle(func(ctx context.Context, input *string) (*string, error) {
			return nil, exceptions.NotFound("resource not found")
		})

		result, err := handler(context.Background(), stringPtr("input"))
		require.Nil(t, result)
		require.Error(t, err)

		// The error should be a huma.StatusError
		var humaErr huma.StatusError
		require.ErrorAs(t, err, &humaErr)
		require.Equal(t, http.StatusNotFound, humaErr.GetStatus())
	})

	t.Run("APIError with details conversion", func(t *testing.T) {
		fields := []huma.ErrorDetail{
			{Message: "invalid format", Location: "body.email"},
		}
		handler := exceptions.Handle(func(ctx context.Context, input *string) (*string, error) {
			return nil, exceptions.Validation("validation failed", fields)
		})

		result, err := handler(context.Background(), stringPtr("input"))
		require.Nil(t, result)
		require.Error(t, err)

		// The error should be a huma.StatusError
		var humaErr huma.StatusError
		require.ErrorAs(t, err, &humaErr)
		require.Equal(t, http.StatusUnprocessableEntity, humaErr.GetStatus())
	})

	t.Run("non-APIError conversion", func(t *testing.T) {
		handler := exceptions.Handle(func(ctx context.Context, input *string) (*string, error) {
			return nil, errors.New("unexpected error")
		})

		result, err := handler(context.Background(), stringPtr("input"))
		require.Nil(t, result)
		require.Error(t, err)

		// Should be converted to 500 error
		var humaErr huma.StatusError
		require.ErrorAs(t, err, &humaErr)
		require.Equal(t, http.StatusInternalServerError, humaErr.GetStatus())
	})

	t.Run("huma StatusError passthrough", func(t *testing.T) {
		humaErr := huma.Error404NotFound("custom not found")
		handler := exceptions.Handle(func(ctx context.Context, input *string) (*string, error) {
			return nil, humaErr
		})

		result, err := handler(context.Background(), stringPtr("input"))
		require.Nil(t, result)
		require.Error(t, err)

		// Should pass through the huma error unchanged
		require.Equal(t, humaErr, err)
	})
}

func stringPtr(s string) *string {
	return &s
}
