package router

import (
	"net/http"

	"github.com/panuwat39/my-api/internal/shared/middleware"
)

func New() http.Handler {
	mux := http.NewServeMux()

	registerV1Routes(mux)

	handler := middleware.RequestID(mux)
	handler = middleware.Logging(handler)
	handler = middleware.Recovery(handler)

	return handler
}
