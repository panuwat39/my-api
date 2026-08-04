package router

import (
	"net/http"

	"github.com/panuwat39/my-api/internal/shared/middleware"
)

func New() http.Handler {
	mux := http.NewServeMux()

	registerV1Routes(mux)

	return middleware.RequestID(mux)
}
