// Package errors provides centralized error handling and common error types
package errors

import (
	"fmt"
	"net/http"
)

// ErrorCode represents different types of application errors
type ErrorCode string

const (
	// Validation errors
	ErrCodeValidation ErrorCode = "VALIDATION_ERROR"
	
	// Not found errors
	ErrCodeNotFound ErrorCode = "NOT_FOUND"
	
	// Storage errors
	ErrCodeStorage ErrorCode = "STORAGE_ERROR"
	
	// Business logic errors
	ErrCodeBusiness ErrorCode = "BUSINESS_ERROR"
	
	// Internal server errors
	ErrCodeInternal ErrorCode = "INTERNAL_ERROR"
)

// AppError represents a structured application error
type AppError struct {
	Code       ErrorCode `json:"code"`
	Message    string    `json:"message"`
	HTTPStatus int       `json:"-"`
	Details    any       `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewValidationError creates a new validation error
func NewValidationError(message string, details any) *AppError {
	return &AppError{
		Code:       ErrCodeValidation,
		Message:    message,
		HTTPStatus: http.StatusBadRequest,
		Details:    details,
	}
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(resource string) *AppError {
	return &AppError{
		Code:       ErrCodeNotFound,
		Message:    fmt.Sprintf("%s não encontrado", resource),
		HTTPStatus: http.StatusNotFound,
	}
}

// NewStorageError creates a new storage error
func NewStorageError(message string, details any) *AppError {
	return &AppError{
		Code:       ErrCodeStorage,
		Message:    message,
		HTTPStatus: http.StatusInternalServerError,
		Details:    details,
	}
}

// NewBusinessError creates a new business logic error
func NewBusinessError(message string, details any) *AppError {
	return &AppError{
		Code:       ErrCodeBusiness,
		Message:    message,
		HTTPStatus: http.StatusUnprocessableEntity,
		Details:    details,
	}
}

// NewInternalError creates a new internal server error
func NewInternalError(message string, details any) *AppError {
	return &AppError{
		Code:       ErrCodeInternal,
		Message:    message,
		HTTPStatus: http.StatusInternalServerError,
		Details:    details,
	}
}

// IsAppError checks if error is an AppError
func IsAppError(err error) (*AppError, bool) {
	if appErr, ok := err.(*AppError); ok {
		return appErr, true
	}
	return nil, false
}
