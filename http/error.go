// github.com/DauletBai/shanraq.org/http/error.go
package http

import "fmt"

// HTTPError represents an error with an associated HTTP status code.
type HTTPError struct {
	Code    int         `json:"-"` // HTTP status code, not marshalled into JSON by default
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"` // Optional additional details
}

// NewHTTPError creates a new HTTPError instance.
func NewHTTPError(code int, message string, details ...interface{}) *HTTPError {
	he := &HTTPError{
		Code:    code,
		Message: message,
	}
	if len(details) > 0 {
		he.Details = details[0] // Assuming first detail is the primary one for now
	}
	return he
}

// Error makes HTTPError satisfy the error interface.
func (e *HTTPError) Error() string {
	if e.Details != nil {
		return fmt.Sprintf("HTTP %d: %s - Details: %v", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("HTTP %d: %s", e.Code, e.Message)
}

// FieldError описывает ошибку для конкретного поля.
type FieldError struct {
	Field   string `json:"field"`   
	Message string `json:"message"` 
}

// ValidationErrors is a list of validation errors.
type ValidationErrors struct {
	Errors []FieldError `json:"validation_errors"` 
}

// NewValidationErrors creates an instance of ValidationErrors.
func NewValidationErrors() *ValidationErrors {
	return &ValidationErrors{
		Errors: make([]FieldError, 0),
	}
}

// Add adds a new field error.
func (ve *ValidationErrors) Add(field, message string) {
	ve.Errors = append(ve.Errors, FieldError{Field: field, Message: message})
}

// IsEmpty checks if there are any validation errors.
func (ve *ValidationErrors) IsEmpty() bool {
	return len(ve.Errors) == 0
}