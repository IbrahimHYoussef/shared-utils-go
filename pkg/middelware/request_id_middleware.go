package middelware

import (
	"context"
	"net/http"

	responses "github.com/IbrahimHYoussef/shared-utils-go/pkg/responses"
	"github.com/google/uuid"
)

// RequestID is the context key type used for storing request IDs.
type RequestID string

// RequestIDKey is the context key used by RequestIDMiddleWare.
const RequestIDKey RequestID = "RequestID"

// RequestIDMiddleWare creates a UUIDv7 request ID, stores it in the request
// context, logs the request, and passes the request to next.
//
// Handlers can read the request ID from r.Context().Value(RequestIDKey).
func RequestIDMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		requestUUID, err := uuid.NewV7()
		if err != nil {
			logger.Error("RequestIDMiddleWare", "error", err)
			responses.RespondWithError(w, http.StatusInternalServerError, "InternalServerError")
			return
		}

		ctx := context.WithValue(r.Context(), RequestIDKey, requestUUID)

		logger.Info("Request", "remote address", r.RemoteAddr, "method", r.Method, "url", r.URL.Path, string(RequestIDKey), requestUUID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
