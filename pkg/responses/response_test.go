package responses_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/IbrahimHYoussef/shared-utils-go/pkg/responses"
)

func TestRespondWithSuccess(t *testing.T) {
	w := httptest.NewRecorder()

	err := responses.RespondWithSuccess(w, http.StatusCreated, map[string]string{"id": "1"})
	if err != nil {
		t.Fatalf("RespondWithSuccess: %s", err)
	}
	if w.Code != http.StatusCreated {
		t.Fatalf("status: got %d, want %d", w.Code, http.StatusCreated)
	}
	if got := w.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("content type: got %q, want application/json", got)
	}
	if got := strings.TrimSpace(w.Body.String()); got != `{"id":"1"}` {
		t.Fatalf("body: got %s", got)
	}
}

func TestRespondWithError(t *testing.T) {
	w := httptest.NewRecorder()

	err := responses.RespondWithError(w, http.StatusUnauthorized, map[string]string{"error": "Token Expired"})
	if err != nil {
		t.Fatalf("RespondWithError: %s", err)
	}
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status: got %d, want %d", w.Code, http.StatusUnauthorized)
	}
	if got := w.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("content type: got %q, want application/json", got)
	}
	if got := strings.TrimSpace(w.Body.String()); got != `{"error":"Token Expired"}` {
		t.Fatalf("body: got %s", got)
	}
}
