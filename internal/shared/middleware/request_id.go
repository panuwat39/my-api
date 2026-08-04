// package middleware

// import (
// 	"crypto/rand"
// 	"encoding/hex"
// 	"net/http"
// )

// const RequestIDHeader = "X-Request-ID"

// func RequestID(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		requestID := r.Header.Get(RequestIDHeader)
// 		if requestID == "" {
// 			requestID = newRequestID()
// 		}

// 		w.Header().Set(RequestIDHeader, requestID)
// 		next.ServeHTTP(w, r)
// 	})
// }

// func newRequestID() string {
// 	buffer := make([]byte, 16)

// 	if _, err := rand.Read(buffer); err != nil {
// 		return "unknown"
// 	}

// 	return hex.EncodeToString(buffer)
// }

package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

const RequestIDHeader = "X-Request-ID"

type contextKey string

const requestIDContextKey contextKey = "request_id"

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(RequestIDHeader)

		if requestID == "" {
			requestID = newRequestID()
		}

		ctx := context.WithValue(
			r.Context(),
			requestIDContextKey,
			requestID,
		)

		w.Header().Set(RequestIDHeader, requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequestIDFromContext(ctx context.Context) string {
	requestID, ok := ctx.Value(requestIDContextKey).(string)
	if !ok {
		return ""
	}

	return requestID
}

func newRequestID() string {
	buffer := make([]byte, 16)

	if _, err := rand.Read(buffer); err != nil {
		return "unknown"
	}

	return hex.EncodeToString(buffer)
}
