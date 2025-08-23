// Package exceptions provides error handling and database driver
// integration for the llm-gateway service.
package exceptions

import (
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type APIError interface {
	error
	Status() int
	Title() string
	Detail() string
	ErrorDetails() []huma.ErrorDetail
}

type BaseAPIError struct {
	status  int
	title   string
	detail  string
	details []huma.ErrorDetail
}

func (e *BaseAPIError) Error() string {
	return fmt.Sprintf("%d %s: %s", e.status, e.title, e.detail)
}

func (e *BaseAPIError) Status() int {
	return e.status
}

func (e *BaseAPIError) Title() string {
	return e.title
}

func (e *BaseAPIError) Detail() string {
	return e.detail
}

func (e *BaseAPIError) ErrorDetails() []huma.ErrorDetail {
	return e.details
}

func New(status int, title, detail string, errs ...huma.ErrorDetail) APIError {
	return &BaseAPIError{
		status:  status,
		title:   title,
		detail:  detail,
		details: errs,
	}
}

func NotFound(detail string) APIError {
	return New(404, "Not Found", detail)
}

func Validation(detail string, fields []huma.ErrorDetail) APIError {
	return New(422, "Validation Failed", detail, fields...)
}

func Unauthorized(detail string) APIError {
	return New(401, "Unauthorized", detail)
}

func Forbidden(detail string) APIError {
	return New(403, "Forbidden", detail)
}

func InternalServerError(detail string) APIError {
	return New(500, "Internal Server Error", detail)
}

// IsNotFound checks if an error is a 404 Not Found APIError.
func IsNotFound(err error) bool {
	apiErr, ok := err.(APIError)
	return ok && apiErr.Status() == http.StatusNotFound
}

// IsUnauthorized checks if an error is a 401 Unauthorized APIError.
func IsUnauthorized(err error) bool {
	apiErr, ok := err.(APIError)
	return ok && apiErr.Status() == http.StatusUnauthorized
}

// IsValidationError checks if an error is a 422 Validation Failed APIError.
func IsValidationError(err error) bool {
	apiErr, ok := err.(APIError)
	return ok && apiErr.Status() == http.StatusUnprocessableEntity
}

// Maybe wraps an error to be directly returned if it matches a condition.
func Maybe(err error, fn func(error) bool, fallback APIError) error {
	if fn(err) {
		return fallback
	}
	return err
}
