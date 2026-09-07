package jsonutil

import (
	"encoding/json"
	"errors"
	"net/http"
)

// DecodeJson decodes the JSON request body into a new T and closes the body.
func DecodeJson[T any](r *http.Request) (*T, error) {
	defer r.Body.Close()
	var req T
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

// ExtractJsonFromCtx is reserved for decoding a JSON value previously stored in
// the request context.
//
// ExtractJsonFromCtx currently returns an unimplemented error.
func ExtractJsonFromCtx[T any](r *http.Request) (*T, error) {
	return nil, errors.New("unimplemented")
}

// EncodeJson writes resBody to w as JSON and sets the Content-Type header to
// application/json.
func EncodeJson[T any](resBody T, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(resBody)
}
