package jsonutil

import (
	"encoding/json"
	"errors"
	"net/http"
)

func DecodeJson[T any](r *http.Request) (*T, error) {
	defer r.Body.Close()
	var req T
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func ExtractJsonFromCtx[T any](r *http.Request) (*T, error) {
	return nil, errors.New("unimplemented")
}

func EncodeJson[T any](resBody T, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(resBody)
}
