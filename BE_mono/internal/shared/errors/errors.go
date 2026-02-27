package errors

import "net/http"

type APIError struct {
	Status  int
	Code    string
	Message string
	Details any
}

func (e *APIError) Error() string {
	return e.Message
}

func New(status int, code string, message string) *APIError {
	return &APIError{Status: status, Code: code, Message: message}
}

var (
	ErrUnauthorized = New(http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
	ErrForbidden    = New(http.StatusForbidden, "FORBIDDEN", "Forbidden")
	ErrNotFound     = New(http.StatusNotFound, "RESOURCE_NOT_FOUND", "Resource not found")
)
