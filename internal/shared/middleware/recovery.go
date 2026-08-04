// package middleware

// import (
// 	"encoding/json"
// 	"log"
// 	"net/http"
// )

// func Recovery(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		defer func() {
// 			if recovered := recover(); recovered != nil {
// 				log.Printf(
// 					"panic recovered: method=%s path=%s error=%v",
// 					r.Method,
// 					r.URL.Path,
// 					recovered,
// 				)

// 				w.Header().Set("Content-Type", "application/json")
// 				w.WriteHeader(http.StatusInternalServerError)

// 				response := map[string]any{
// 					"error": map[string]string{
// 						"code":    "INTERNAL_SERVER_ERROR",
// 						"message": "internal server error",
// 					},
// 				}

// 				if err := json.NewEncoder(w).Encode(response); err != nil {
// 					log.Printf("write recovery response error: %v", err)
// 				}
// 			}
// 		}()

// 		next.ServeHTTP(w, r)
// 	})
// }

package middleware

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"runtime/debug"
)

func Recovery(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			recovered := recover()
			if recovered == nil {
				return
			}

			logger.ErrorContext(
				r.Context(),
				"panic recovered",
				slog.String(
					"request_id",
					RequestIDFromContext(r.Context()),
				),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Any("error", recovered),
				slog.String("stack", string(debug.Stack())),
			)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)

			response := map[string]any{
				"error": map[string]string{
					"code":    "INTERNAL_SERVER_ERROR",
					"message": "internal server error",
				},
			}

			if err := json.NewEncoder(w).Encode(response); err != nil {
				logger.ErrorContext(
					r.Context(),
					"write recovery response failed",
					slog.Any("error", err),
				)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
