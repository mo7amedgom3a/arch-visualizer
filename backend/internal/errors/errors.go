package errors

import (
	"fmt"
	"strings"
)

// ErrorKind represents the category of error
type ErrorKind string

const (
	// KindValidation indicates a validation error
	KindValidation ErrorKind = "validation"
	// KindNotFound indicates a resource was not found
	KindNotFound ErrorKind = "not_found"
	// KindConflict indicates a conflict (e.g., duplicate resource)
	KindConflict ErrorKind = "conflict"
	// KindUnauthorized indicates an authorization error
	KindUnauthorized ErrorKind = "unauthorized"
	// KindForbidden indicates a forbidden operation
	KindForbidden ErrorKind = "forbidden"
	// KindInternal indicates an internal server error
	KindInternal ErrorKind = "internal"
	// KindBadRequest indicates a bad request
	KindBadRequest ErrorKind = "bad_request"
	// KindTimeout indicates a timeout error
	KindTimeout ErrorKind = "timeout"
	// KindUnavailable indicates a service is unavailable
	KindUnavailable ErrorKind = "unavailable"
)

// AppError represents a standardized application error
type AppError struct {
	// Code is a unique error code (e.g., "RESOURCE_NOT_FOUND", "INVALID_INPUT")
	Code string `json:"code"`
	// Kind is the error category
	Kind ErrorKind `json:"kind"`
	// Message is a human-readable error message
	Message string `json:"message"`
	// Op is the operation where the error occurred (e.g., "ComputeService.CreateInstance")
	Op string `json:"op,omitempty"`
	// Cause is the underlying error that caused this error
	Cause error `json:"cause,omitempty"`
	// Metadata contains additional context about the error
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e == nil {
		return ""
	}

	var parts []string
	if e.Op != "" {
		parts = append(parts, e.Op)
	}
	if e.Code != "" {
		parts = append(parts, fmt.Sprintf("[%s]", e.Code))
	}
	if e.Message != "" {
		parts = append(parts, e.Message)
	}
	if e.Cause != nil {
		parts = append(parts, fmt.Sprintf(": %v", e.Cause))
	}

	return strings.Join(parts, " ")
}

// Unwrap returns the underlying error for error wrapping compatibility
func (e *AppError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

// New creates a new AppError with the given code, kind, and message
func New(code string, kind ErrorKind, message string) *AppError {
	return &AppError{
		Code:    code,
		Kind:    kind,
		Message: message,
	}
}

// Newf creates a new AppError with formatted message
func Newf(code string, kind ErrorKind, format string, args ...interface{}) *AppError {
	return &AppError{
		Code:    code,
		Kind:    kind,
		Message: fmt.Sprintf(format, args...),
	}
}

// Wrap wraps an existing error into an AppError
func Wrap(err error, code string, kind ErrorKind, message string) *AppError {
	if err == nil {
		return nil
	}

	// If err is already an AppError, preserve its metadata and add new context
	if appErr, ok := err.(*AppError); ok {
		return &AppError{
			Code:     code,
			Kind:     kind,
			Message:  message,
			Cause:    appErr,
			Metadata: appErr.Metadata,
		}
	}

	return &AppError{
		Code:    code,
		Kind:    kind,
		Message: message,
		Cause:   err,
	}
}

// Wrapf wraps an existing error with a formatted message
func Wrapf(err error, code string, kind ErrorKind, format string, args ...interface{}) *AppError {
	if err == nil {
		return nil
	}

	// If err is already an AppError, preserve its metadata and add new context
	if appErr, ok := err.(*AppError); ok {
		return &AppError{
			Code:     code,
			Kind:     kind,
			Message:  fmt.Sprintf(format, args...),
			Cause:    appErr,
			Metadata: appErr.Metadata,
		}
	}

	return &AppError{
		Code:    code,
		Kind:    kind,
		Message: fmt.Sprintf(format, args...),
		Cause:   err,
	}
}

// WithOp sets the operation context for the error
func (e *AppError) WithOp(op string) *AppError {
	if e == nil {
		return nil
	}
	e.Op = op
	return e
}

// WithMeta adds metadata to the error
func (e *AppError) WithMeta(key string, value interface{}) *AppError {
	if e == nil {
		return nil
	}
	if e.Metadata == nil {
		e.Metadata = make(map[string]interface{})
	}
	e.Metadata[key] = value
	return e
}

// WithMetadata sets the entire metadata map
func (e *AppError) WithMetadata(metadata map[string]interface{}) *AppError {
	if e == nil {
		return nil
	}
	e.Metadata = metadata
	return e
}

// IsKind checks if the error is of the given kind
func IsKind(err error, kind ErrorKind) bool {
	if err == nil {
		return false
	}

	appErr, ok := err.(*AppError)
	if !ok {
		return false
	}

	return appErr.Kind == kind
}

// IsCode checks if the error has the given code
func IsCode(err error, code string) bool {
	if err == nil {
		return false
	}

	appErr, ok := err.(*AppError)
	if !ok {
		return false
	}

	return appErr.Code == code
}

// AsAppError extracts the AppError from an error chain
func AsAppError(err error) *AppError {
	if err == nil {
		return nil
	}

	appErr, ok := err.(*AppError)
	if ok {
		return appErr
	}

	// Try to unwrap and check again
	if unwrapped := Unwrap(err); unwrapped != nil {
		return AsAppError(unwrapped)
	}

	return nil
}

// Unwrap is a helper that works with standard error wrapping
func Unwrap(err error) error {
	if err == nil {
		return nil
	}

	type unwrapper interface {
		Unwrap() error
	}

	if u, ok := err.(unwrapper); ok {
		return u.Unwrap()
	}

	return nil
}
