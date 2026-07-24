package errorsutil

import (
	"net/http"

	responses "github.com/IbrahimHYoussef/shared-utils-go/pkg/responses"
)

// ServiceError represents an error with an HTTP status code
type ServiceError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"error"`
}

// Error implements the error interface
func (e *ServiceError) Error() string {
	return e.Message
}

func NewError(status_code int, message string) *ServiceError {
	return &ServiceError{StatusCode: status_code, Message: message}
}

// Helper functions to create common errors
func NewBadRequestError(message string) *ServiceError {
	return &ServiceError{StatusCode: http.StatusBadRequest, Message: message}
}

func NewUnauthorizedError(message string) *ServiceError {
	return &ServiceError{StatusCode: http.StatusUnauthorized, Message: message}
}

func NewForbiddenError(message string) *ServiceError {
	return &ServiceError{StatusCode: http.StatusForbidden, Message: message}
}

func NewNotFoundError(message string) *ServiceError {
	return &ServiceError{StatusCode: http.StatusNotFound, Message: message}
}

func NewConflictError(message string) *ServiceError {
	return &ServiceError{StatusCode: http.StatusConflict, Message: message}
}

func NewInternalError(message string) *ServiceError {
	return &ServiceError{StatusCode: http.StatusInternalServerError, Message: message}
}

func HandleServiceError(w http.ResponseWriter, err error) {
	if serviceErr, ok := err.(*ServiceError); ok {
		responses.RespondWithError(w, serviceErr.StatusCode, serviceErr)
	} else {
		responses.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
	}
}
