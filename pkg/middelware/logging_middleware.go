package middelware

import (
	"net/http"
)

func LoggingMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Request", "remote address", r.RemoteAddr, "method", r.Method, "url", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
