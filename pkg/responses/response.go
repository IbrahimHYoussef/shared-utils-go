package responses

import (
	"net/http"

	"github.com/IbrahimHYoussef/shared-utils-go/pkg/jsonutil"
)

type UnifiedResponse[T any, E any] struct {
	StatusCode int `json:"status_code"`
	Body       T   `json:"body,omitempty"`
	Error      E   `json:"error,omitempty"`
	isError    bool
}

func New[T any, E any](statusCode int, body T, error E, isError bool) UnifiedResponse[T, E] {
	return UnifiedResponse[T, E]{
		StatusCode: statusCode,
		Body:       body,
		Error:      error,
		isError:    isError,
	}
}

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

func RespondWithSuccess[T any](w http.ResponseWriter, statusCode int, successResponse T) error {
	response := New[T, any](statusCode, successResponse, nil, false)
	return response.ReturnResponse(w)
}

func RespondWithError[E any](w http.ResponseWriter, statusCode int, errorResponse E) error {
	response := New[any, E](statusCode, nil, errorResponse, true)
	return response.ReturnResponse(w)
}
