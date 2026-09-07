// Package responses contains helpers for writing consistent JSON HTTP
// responses.
package responses

import (
	"net/http"

	"github.com/IbrahimHYoussef/shared-utils-go/pkg/jsonutil"
)

// UnifiedResponse represents either a successful response body or an error body
// with an HTTP status code.
type UnifiedResponse[T any, E any] struct {
	// StatusCode is the HTTP status code to write.
	StatusCode int `json:"status_code"`
	// Body is the success response payload.
	Body T `json:"body,omitempty"`
	// Error is the error response payload.
	Error   E `json:"error,omitempty"`
	isError bool
}

// New returns a UnifiedResponse with the supplied status, body, error, and mode.
//
// When isError is true, ReturnResponse writes error. Otherwise it writes body.
func New[T any, E any](statusCode int, body T, error E, isError bool) UnifiedResponse[T, E] {
	return UnifiedResponse[T, E]{
		StatusCode: statusCode,
		Body:       body,
		Error:      error,
		isError:    isError,
	}
}

// ReturnResponse writes the configured HTTP status code and JSON payload to w.
func (ur UnifiedResponse[T, E]) ReturnResponse(w http.ResponseWriter) error {
	w.WriteHeader(ur.StatusCode)

	var body any
	if ur.isError == false {
		body = ur.Body
	} else {
		body = ur.Error
	}

	err := jsonutil.EncodeJson(body, w)
	if err != nil {
		return err
	}
	return nil
}

// RespondWithSuccess writes successResponse as a JSON response with statusCode.
func RespondWithSuccess[T any](w http.ResponseWriter, statusCode int, successResponse T) error {
	response := New[T, any](statusCode, successResponse, nil, false)
	return response.ReturnResponse(w)
}

// RespondWithError writes errorResponse as a JSON error response with statusCode.
func RespondWithError[E any](w http.ResponseWriter, statusCode int, errorResponse E) error {
	response := New[any, E](statusCode, nil, errorResponse, true)
	return response.ReturnResponse(w)
}
