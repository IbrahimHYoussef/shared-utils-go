// Package errorsutil provides HTTP-aware error helpers for service and handler
// layers.
package errorsutil

import (
	"net/http"

	responses "github.com/IbrahimHYoussef/shared-utils-go/pkg/responses"
)

// ServiceError represents an error with an HTTP status code.
type ServiceError struct {
	// StatusCode is the HTTP status code that should be sent to the client.
	StatusCode int `json:"status_code"`
	// Message is the user-facing error message.
	Message string `json:"error"`
}

// Error implements the error interface.
func (e *ServiceError) Error() string {
	return e.Message
}

// NewError returns a ServiceError with the supplied HTTP status code and
// message.
func NewError(status_code int, message string) *ServiceError {
	return &ServiceError{StatusCode: status_code, Message: message}
}

// NewBadRequestError returns a ServiceError with status 400 Bad Request.
func NewBadRequestError(message string) *ServiceError {
	return &ServiceError{StatusCode: http.StatusBadRequest, Message: message}
}

// NewUnauthorizedError returns a ServiceError with status 401 Unauthorized.
func NewUnauthorizedError(message string) *ServiceError {
	return &ServiceError{StatusCode: http.StatusUnauthorized, Message: message}
}

// NewForbiddenError returns a ServiceError with status 403 Forbidden.
func NewForbiddenError(message string) *ServiceError {
	return &ServiceError{StatusCode: http.StatusForbidden, Message: message}
}

// NewNotFoundError returns a ServiceError with status 404 Not Found.
func NewNotFoundError(message string) *ServiceError {
	return &ServiceError{StatusCode: http.StatusNotFound, Message: message}
}

// NewConflictError returns a ServiceError with status 409 Conflict.
func NewConflictError(message string) *ServiceError {
	return &ServiceError{StatusCode: http.StatusConflict, Message: message}
}

// NewInternalError returns a ServiceError with status 500 Internal Server Error.
func NewInternalError(message string) *ServiceError {
	return &ServiceError{StatusCode: http.StatusInternalServerError, Message: message}
}

// HandleServiceError writes err to w as a JSON error response.
//
// If err is a *ServiceError, its status code and body are used. Any other error
// is returned as a generic 500 Internal Server Error response.
func HandleServiceError(w http.ResponseWriter, err error) {
	if serviceErr, ok := err.(*ServiceError); ok {
		responses.RespondWithError(w, serviceErr.StatusCode, serviceErr)
	} else {
		responses.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
	}
}
