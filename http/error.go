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