package exceptions

import (
	"context"

	"github.com/WebDeveloperBen/ai-gateway/internal/logger"
	"github.com/danielgtaylor/huma/v2"
)

type HandlerFunc[In any, Out any] func(ctx context.Context, input *In) (*Out, error)

// Handle is a helper to wrap Huma-compatible handlers with centralized error handling.
//
// It allows your handlers to return structured domain-specific errors (APIError),
// which are converted into `huma` error responses automatically.
//
// If the handler returns a known `APIError`, the wrapper will:
//   - Log the error using the custom logger (`lib.LogError`)
//   - Convert it to a `huma.NewError(...)` with appropriate status and detail
//
// If the error is unknown or not an `APIError`, it returns a generic 500 error.
//
// Example usage:
//
//	huma.Get(api, "/users/{id}", exceptions.Wrap(func(ctx context.Context, input *Input) (*Output, error) {
//	    if input.ID == "123" {
//	        return nil, exceptions.NotFound("user not found")
//	    }
//	    return &Output{...}, nil
//	}))
//
// This lets you keep your business logic clean while handling errors in a consistent,
// OpenAPI-friendly format.
func Handle[In any, Out any](handler HandlerFunc[In, Out]) HandlerFunc[In, Out] {
	return func(ctx context.Context, input *In) (*Out, error) {
		out, err := handler(ctx, input)
		if err == nil {
			return out, nil
		}

		logger.LogError(ctx, err, "handler error")

		if apiErr, ok := err.(APIError); ok {
			detailErrors := make([]error, 0, len(apiErr.ErrorDetails()))
			for _, d := range apiErr.ErrorDetails() {
				detailErrors = append(detailErrors, &d)
			}
			return nil, huma.NewError(apiErr.Status(), apiErr.Detail(), detailErrors...)
		}

		// If it's already a huma StatusError, pass it through unchanged
		if _, ok := err.(huma.StatusError); ok {
			return nil, err
		}

		return nil, huma.Error500InternalServerError("unexpected server error")
	}
}
