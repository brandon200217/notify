package middleware

import (
	"net/http"

	"github.com/brandon200217/NOTIFY/internal/logger"
	"github.com/google/uuid"
)

const RequestIDHeader = "X-Request-ID"

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get(RequestIDHeader)
		if reqID == "" {
			reqID = uuid.NewString()
		}

		ctx := logger.WithRequestID(r.Context(), reqID)
		r = r.WithContext(ctx)

		w.Header().Set(RequestIDHeader, reqID)

		next.ServeHTTP(w, r)
	})
}
