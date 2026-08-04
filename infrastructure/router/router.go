// package router

// import (
// 	"net/http"

// 	"github.com/panuwat39/my-api/internal/shared/middleware"
// )

// func New() http.Handler {
// 	mux := http.NewServeMux()

// 	registerV1Routes(mux)

// 	handler := middleware.RequestID(mux)
// 	handler = middleware.Logging(handler)
// 	handler = middleware.Recovery(handler)

// 	return handler
// }

package router

import (
	"log/slog"
	"net/http"

	"github.com/panuwat39/my-api/internal/shared/middleware"
)

func New(logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	registerV1Routes(mux)

	var handler http.Handler = mux

	handler = middleware.Logging(logger, handler)
	handler = middleware.RequestID(handler)
	handler = middleware.Recovery(logger, handler)

	return handler
}
