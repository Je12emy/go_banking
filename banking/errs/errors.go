package errs

import "net/http"

type AppError struct {
	Code    int `json:",omitempty"`
	Message string
}

// AsMessage : Returns the error message only, Code will be empty
func (e AppError) AsMessage() *AppError {
	return &AppError{
		Message: e.Message,
	}
}

// NewNotFoundError : Returns a not found error message
func NewNotFoundError(message string) *AppError {
	return &AppError{
		Message: message,
		Code:    http.StatusNotFound,
	}
}

// NewUnexpectedError : Returns a Internal Server Error
func NewUnexpectedError(message string) *AppError {
	return &AppError{
		Message: message,
		Code:    http.StatusInternalServerError,
	}
}

// NewValidationError : Returns a validation error
func NewValidationError(message string) *AppError {
	// StatusUnprocessableEntity == Bussiness rules deny this request
	return &AppError{
		Message: message,
		Code:    http.StatusUnprocessableEntity,
	}
}
