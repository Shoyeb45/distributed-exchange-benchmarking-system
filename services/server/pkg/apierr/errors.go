package apierr

import (
	"errors"
	"fmt"
	"net/http"
)

type ErrorCode string

const (
	CodeBadRequest          ErrorCode = "BAD_REQUEST"
	CodeUnauthorized        ErrorCode = "UNAUTHORIZED"
	CodeForbidden           ErrorCode = "FORBIDDEN"
	CodeNotFound            ErrorCode = "NOT_FOUND"
	CodeConflict            ErrorCode = "CONFLICT"
	CodeUnprocessableEntity ErrorCode = "UNPROCESSABLE_ENTITY"
	CodeTooManyRequests     ErrorCode = "TOO_MANY_REQUESTS"
	CodeInternal            ErrorCode = "INTERNAL_SERVER_ERROR"
	CodeServiceUnavailable  ErrorCode = "SERVICE_UNAVAILABLE"
)

type APIError struct {
	// http status code
	StatusCode int `json:"-"`
	// readable error by client to handle the error
	Code ErrorCode `json:"code"`
	// user level message
	Message string `json:"message"`
	// field level validation errors
	Details []FieldError `json:"details,omitempty"`
	// Internal cause, logged and never exposed in prod
	cause error
}

// NewBadRequest - 400.
func NewBadRequest(message string) *APIError {
	if message == "" {
		message = "bad request"
	}
	return &APIError{StatusCode: http.StatusBadRequest, Code: CodeBadRequest, Message: message}
}

// NewUnauthorized - 401, using for unavailability of auth.
func NewUnauthorized(message string) *APIError {
	if message == "" {
		message = "authentication required"
	}
	return &APIError{StatusCode: http.StatusUnauthorized, Code: CodeUnauthorized, Message: message}
}

// NewForbidden — 403, user is authenticated but lacks permission.
func NewForbidden(message string) *APIError {
	if message == "" {
		message = "you do not have permission to perform this action"
	}
	return &APIError{StatusCode: http.StatusForbidden, Code: CodeForbidden, Message: message}
}

// NewNotFound — 404, use when a requested resource does not exist.
func NewNotFound(resource string) *APIError {
	return &APIError{
		StatusCode: http.StatusNotFound,
		Code:       CodeNotFound,
		Message:    fmt.Sprintf("%s not found", resource),
	}
}

// NewMethodNotAllowed — 405, use when a requested api method is not allowed.
func NewMethodNotAllowed(resource string) *APIError {
	return &APIError{
		StatusCode: http.StatusMethodNotAllowed,
		Code:       CodeNotFound,
		Message:    resource,
	}
}

// NewConflict — 409, use for duplicate resource.
func NewConflict(message string) *APIError {
	if message == "" {
		message = "conflict with existing resource"
	}
	return &APIError{StatusCode: http.StatusConflict, Code: CodeConflict, Message: message}
}

// NewValidationError — 422, use for semantic validation failures with field details.
func NewValidationError(message string, details ...FieldError) *APIError {
	return &APIError{
		StatusCode: http.StatusUnprocessableEntity,
		Code:       CodeUnprocessableEntity,
		Message:    message,
		Details:    details,
	}
}

// NewTooManyRequests — 429, Use for rate limiting.
func NewTooManyRequests(message string) *APIError {
	if message == "" {
		message = "too many requests, please slow down"
	}
	return &APIError{StatusCode: http.StatusTooManyRequests, Code: CodeTooManyRequests, Message: message}
}

// NewInternal — 500, Use for unexpected errors. Always attach the real cause via `WithCause`.
func NewInternal(cause error) *APIError {
	return &APIError{
		StatusCode: http.StatusInternalServerError,
		Code:       CodeInternal,
		Message:    "an unexpected error occurred",
		cause:      cause,
	}
}

// NewServiceUnavailable — 503, use when a downstream dependency (DB, cache, etc.) is down.
func NewServiceUnavailable(message string) *APIError {
	if message == "" {
		message = "service temporarily unavailable"
	}
	return &APIError{StatusCode: http.StatusServiceUnavailable, Code: CodeServiceUnavailable, Message: message}
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *APIError) Unwrap() error {
	return e.cause
}

func (e *APIError) WithCause(err error) *APIError {
	e.cause = err
	return e
}

func (e *APIError) WithDetails(details ...FieldError) *APIError {
	e.Details = append(e.Details, details...)
	return e
}

func (e *APIError) Cause() error {
	return e.cause
}

func As(err error, target **APIError) bool {
	return errors.As(err, target)
}

func IsAPIError(err error) (*APIError, bool) {
	var apierr *APIError

	for e := err; e != nil; {
		ae := &APIError{}
		if errors.As(e, &ae) {
			return ae, true
		}
		type unwrapper interface{ Unwrap() error }

		if u, ok := e.(unwrapper); ok {
			e = u.Unwrap()
		} else {
			break
		}
	}
	_ = apierr

	return nil, false
}
