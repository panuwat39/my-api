// package middleware

// import (
// 	"log"
// 	"net/http"
// 	"time"
// )

// type responseWriter struct {
// 	http.ResponseWriter
// 	status int
// }

// func (w *responseWriter) WriteHeader(statusCode int) {
// 	w.status = statusCode
// 	w.ResponseWriter.WriteHeader(statusCode)
// }

// func Logging(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		start := time.Now()

// 		writer := &responseWriter{
// 			ResponseWriter: w,
// 			status:         http.StatusOK,
// 		}

// 		next.ServeHTTP(writer, r)

// 		log.Printf(
// 			"request_id=%s method=%s path=%s status=%d duration=%s",
// 			r.Header.Get(RequestIDHeader),
// 			r.Method,
// 			r.URL.Path,
// 			writer.status,
// 			time.Since(start),
// 		)
// 	})
// }

package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func Logging(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		writer := &responseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		next.ServeHTTP(writer, r)

		logger.InfoContext(
			r.Context(),
			"http request completed",
			slog.String(
				"request_id",
				RequestIDFromContext(r.Context()),
			),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", writer.status),
			slog.Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			),
		)
	})
}
