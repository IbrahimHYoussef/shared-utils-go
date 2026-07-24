package middelware

import (
	"context"
	"net/http"

	responses "github.com/IbrahimHYoussef/project-management/shared-utils-go/Responses"
	"github.com/google/uuid"
)

type RequestID string

const RequestIDKey RequestID = "RequestID"

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
